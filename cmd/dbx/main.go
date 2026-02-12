package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
	"github.com/iancoleman/strcase"
)

//go:embed dbx.template
var tmpFS embed.FS

type spec []string

func (t spec) String() string {
	return strings.Join(t, "")
}

func (t spec) IsPointer() bool {
	return len(t) > 0 && t[0] == "*"
}

func (t spec) IsInterface() bool {
	return len(t) > 0 && t[0] == "interface"
}

func (t spec) IsSlice() bool {
	return len(t) > 1 && t[0] == "[" && t[1] == "]"
}

func (t spec) IsMap() bool {
	return len(t) > 0 && t[0] == "map["
}

type field struct {
	Name  string
	Alias string
	Type  spec
}

func getTag(tag, key, not string) *structtag.Tag {
	tt, err := structtag.Parse(tag)
	if err != nil {
		return nil
	}
	t, err := tt.Get(key)
	if err != nil || t.HasOption(not) || t.Name == not || t.Name == "" || t.Name == "-" {
		return nil
	}
	return t
}

func getType(x ast.Expr) []string {
	switch t := x.(type) {
	case *ast.MapType:
		return append(append(append([]string{"map["}, getType(t.Key)...), "]"), getType(t.Value)...)
	case *ast.ArrayType:
		return append(append(append([]string{"["}, getType(t.Len)...), "]"), getType(t.Elt)...)
	case *ast.SelectorExpr:
		return append(append(getType(t.X), "."), t.Sel.Name)
	case *ast.StarExpr:
		return append([]string{"*"}, getType(t.X)...)
	case *ast.BasicLit:
		return []string{t.Value}
	case *ast.Ident:
		return []string{t.Name}
	case *ast.InterfaceType:
		return []string{"interface"}
	case nil:
		return nil
	default:
		panic(t)
	}
}

func getFields(pkg *ast.Package, a []*ast.Field, key, not string) []field {
	var ff []field
	for _, f := range a {
		if f.Tag == nil {
			if f.Names == nil {
				ff = append(ff, getStruct(pkg, f.Type, key, not)...)
			}
			continue
		}
		tg := getTag(strings.Trim(f.Tag.Value, "` "), key, not)
		if tg == nil {
			continue
		}
		if len(f.Names) != 1 {
			panic(f)
		}
		fn := f.Names[0]
		tt := getType(f.Type)
		if tt == nil {
			panic(f)
		}
		ef := field{
			Name:  fn.Name,
			Alias: tg.Name,
			Type:  tt,
		}
		ff = append(ff, ef)
	}
	return ff
}

func getStruct(pkg *ast.Package, e ast.Expr, key, not string) []field {
	switch t := e.(type) {
	case *ast.StructType:
		return getFields(pkg, t.Fields.List, key, not)
	case *ast.Ident:
		return getStructFields(pkg, t.Name, key, not)
	case *ast.StarExpr:
		return getStruct(pkg, t.X, key, not)
	default:
		panic(t)
	}
}

func getStructFields(pkg *ast.Package, name, key, not string) []field {
	for _, file := range pkg.Files {
		for _, dl := range file.Decls {
			if gd, ok := dl.(*ast.GenDecl); ok && gd != nil && gd.Tok == token.TYPE {
				for _, sp := range gd.Specs {
					switch t := sp.(type) {
					case *ast.TypeSpec:
						if t.Name.Name == name {
							return getStruct(pkg, t.Type, key, not)
						}
					default:
						panic(t)
					}
				}
			}
		}
	}
	return nil
}

type class struct {
	Type   string
	Fields []field
}

func (t class) Last() int {
	return len(t.Fields) - 1
}

func getTypes(pkg *ast.Package, key, not string, names ...string) []class {
	var types []class
	for _, name := range names {
		fields := getStructFields(pkg, name, key, not)
		if fields != nil {
			types = append(types, class{Type: name, Fields: uniqueFields(fields)})
		}
	}
	return types
}

func uniqueFields(fields []field) []field {
	names := map[string]int{}
	for i, n := 0, len(fields); i < n; i++ {
		if j, ok := names[fields[i].Name]; ok {
			copy(fields[j:], fields[j+1:])
			i--
			n--
		}
		names[fields[i].Name] = i
	}
	return fields[:len(names)]
}

type files map[string]struct{}

func (f files) add(name string) {
	f[name] = struct{}{}
}

func (f files) match(name string) bool {
	if path.Ext(name) != ".go" {
		return false
	}
	_, ok := f[name[:len(name)-3]]
	return ok
}

type generate struct {
	Package string
	Types   []class
}

func main() {
	var output, suffix, direct string
	var prefix bool
	flag.StringVar(&suffix, "n", "gen.go", "output file name")
	flag.StringVar(&direct, "w", ".", "work directory")
	flag.BoolVar(&prefix, "x", false, "without prefix")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	filter := files{}
	for _, name := range flag.Args() {
		name = strcase.ToSnake(name)
		if output == "" {
			if !prefix {
				output = name + "_"
			}
			output += suffix
		}
		filter.add(name)
	}
	dir, err := parser.ParseDir(token.NewFileSet(), direct, func(info fs.FileInfo) bool {
		return filter.match(info.Name())
	}, 0)
	if err != nil {
		return
	}

	var types []class
	for _, pkg := range dir {
		t := getTypes(pkg, "db", "password", flag.Args()...)
		if t != nil {
			types = append(types, t...)
		}
	}
	sort.SliceStable(types, func(i, j int) bool {
		return types[i].Type < types[j].Type
	})

	fmt.Printf("%+#v\n", types)

	for _, name := range flag.Args() {
		i := sort.Search(len(types), func(i int) bool {
			return types[i].Type >= name
		})
		if i >= len(types) || types[i].Type != name {
			log.Printf("type `%s` not found in source files", name)
		}
	}

	tmp, err := template.ParseFS(tmpFS, "*")
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(path.Join(direct, output))
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	for name := range dir {
		err = tmp.Execute(&buf, generate{
			Package: name,
			Types:   types,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = buf.WriteTo(out)
	if err != nil {
		return
	}
	err = out.Sync()
	if err != nil {
		log.Fatal(err)
	}
	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}
