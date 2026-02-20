package builder

import (
	"fmt"
	"github.com/pshvedko/dbx/filter"
)

type Key struct {
	Builder
}

func (k *Key) AppendValue(t filter.Type, v any) (int, error) {
	var n1, n2 int
	var err error
	switch t {
	case filter.TRUE, filter.FALSE:
	default:
		n1, err = k.WriteString(t.String())
		if err != nil {
			return n1, err
		}
	}
	z := k.Size()
	f := k.Value(v)
	switch x := f.(type) {
	case Keyword:
		return n1 + 1, k.WriteByte(x[0])
	default:
		if x != nil && z == k.Size() {
			n2, err = k.AppendFormat(x)
			return n1 + n2, err
		}
		return n1, nil
	}
}

func (k *Key) AppendColumn(t string, c string, m map[string]int) (int, error) {
	return k.AppendInt(m[c])
}

func (k *Key) AppendVia(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('&')
	}
	return 1, k.WriteByte('|')
}

func (k *Key) AppendParenthesis(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('[')
	}
	return 1, k.WriteByte(']')
}

func (k *Key) Straight(j filter.Projector, u bool, x filter.Filter) (err error) {
	return Straight(k, j, u, x)
}

func (k *Key) Conjunct(j filter.Projector, u bool, ff []filter.Filter) (err error) {
	_, err = k.AppendParenthesis(true)
	if err != nil {
		return
	}
	for i, f := range ff {
		if i > 0 {
			_, err = k.AppendVia(u)
			if err != nil {
				return
			}
		}
		err = f.To(k, j)
		if err != nil {
			return
		}
	}
	_, err = k.AppendParenthesis(false)
	return
}

type Origin struct {
	File string
	Line int
}

type Hash struct {
	Origin
	static string
}

type Keeper interface {
	Get(Hash) (string, bool)
	Put(Hash, string)
}

type Cache struct {
	Origin
	Keeper
}

func (c Cache) IsCacheEnabled() bool {
	return c.Keeper != nil && c.Origin.File != "" && c.Origin.Line != 0
}

func WritePresentRanges(idx []int) {
	y := true // разрешено печатать индексы
	z := 0    // индекс начала интервала
	for i, x := range idx {
		if x == 0 { // поле не в списке, пропускаем
			if !y { // было запрещено печатать,
				y = true // разрешаем печатать
				switch i - z {
				case 0: // такое невозможно!
					panic(i)
				case 1: // напечатан только один, добавить нечего
				case 2: // напечатано два, между соседними "-" не влезет
					println(",")
					println(i - 1) // конец интервала
				default: // больше двух, ставим "-" закрывая интервал
					println("-")
					println(i - 1) // конец интервала
				}
				println(",") // закрываем интервал
			}
			// тут было разрешено печатать, но нечего
			continue
		}

		if y { // можно печатать индексы
			println(i) // печатаем начало интервала
			z = i      // запоминаем начало интервала
			y = false  // запретим печатать индексы
		}
		// тут нельзя печатать индексы
	}

	if !y { // вышли из цикла и запрещено печатать
		println("-")
	}

	println("===", fmt.Sprint(idx))
}
