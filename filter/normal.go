package filter

type True struct{}

func (f True) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f True) To(b Builder, _ Projector) error {
	_, err := b.AppendValue(TRUE, true)
	return err
}

func (f True) Type() Type { return TRUE }

type False struct{}

func (f False) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f False) To(b Builder, _ Projector) error {
	_, err := b.AppendValue(FALSE, false)
	return err
}

func (f False) Type() Type { return FALSE }

func Expand(t Filter) Filter {
	switch f := t.(type) {
	case Eq:
		a := make(And, 0, len(f))
		for k, v := range f {
			a = append(a, Eq{k: v})
		}
		return a
	case Ne:
		a := make(And, 0, len(f))
		for k, v := range f {
			a = append(a, Ne{k: v})
		}
		return a
	case In:
		o := make(Or, 0, 2)
		for k, x := range f {
			for _, v := range x {
				o = append(o, Eq{k: v})
			}
		}
		if len(o) > 0 {
			return o
		}
		return False{}
	case Ni:
		a := make(And, 0, 2)
		for k, x := range f {
			for _, v := range x {
				a = append(a, Ne{k: v})
			}
		}
		if len(a) > 0 {
			return a
		}
		return True{}
	case And:
		var a And
		var oo []Or
		for _, v := range f {
			switch x := Expand(v).(type) {
			case And:
				a = append(a, x...)
			case Or:
				oo = append(oo, x)
			default:
				a = append(a, x)
			}
		}
		var o Or
		if len(oo) > 0 {
			for _, x := range oo[0] {
				y := a
				y = append(y, x)
				for _, z := range oo[1:] {
					y = append(y, z)
				}
				o = append(o, Expand(y))
			}
		}
		if len(o) > 0 {
			return Expand(o)
		}
		return a
	case Or:
		var o Or
		for _, v := range f {
			switch x := Expand(v).(type) {
			case Or:
				o = append(o, x...)
			default:
				o = append(o, x)
			}
		}
		return o
	default:
		return f
	}
}

type Type int

func (t Type) String() string {
	return TypeName[t]
}

var TypeName = []string{
	EQ:    "EQ",
	NE:    "NE",
	GE:    "GE",
	GT:    "GT",
	LE:    "LE",
	LT:    "LT",
	AS:    "AS",
	NA:    "NA",
	IN:    "IN",
	NI:    "NI",
	FALSE: "FALSE",
	TRUE:  "TRUE",
	AND:   "AND",
	OR:    "OR",
}

const (
	EQ Type = iota
	NE
	GE
	GT
	LE
	LT
	AS
	NA
	IN
	NI
	FALSE
	TRUE
	AND
	OR
)

func Merge[M interface {
	~map[string]V
	Filter
}, V any](mm []M) []M {
	switch len(mm) {
	case 0:
		return nil
	case 1:
		return mm[:1]
	default:
		var n int
		for i := range mm[1:] {
			var z bool
			for _, m := range mm[i+1:] {
				z = false
				for k, v := range mm[i] {
					if _, ok := m[k]; !ok {
						delete(mm[i], k)
						m[k] = v
						continue
					}
					z = true
				}
			}
			if !z {
				n++
			}
		}
		return mm[n:]
	}
}

func Unite(a And, eqs []Eq, nes []Ne, ges []Ge, gts []Gt, les []Le, lts []Lt, ass []As, nas []Na, ins []In, nis []Ni, t, f int) Filter {
	if len(a) == 0 {
		for _, m := range Merge(eqs) {
			a = append(a, m)
		}
		for _, m := range Merge(nes) {
			a = append(a, m)
		}
		for _, m := range Merge(ges) {
			a = append(a, m)
		}
		for _, m := range Merge(gts) {
			a = append(a, m)
		}
		for _, m := range Merge(les) {
			a = append(a, m)
		}
		for _, m := range Merge(lts) {
			a = append(a, m)
		}
		for _, m := range Merge(ass) {
			a = append(a, m)
		}
		for _, m := range Merge(nas) {
			a = append(a, m)
		}
		for _, m := range Merge(ins) {
			a = append(a, m)
		}
		for _, m := range Merge(nis) {
			a = append(a, m)
		}
		if t > 0 {
			a = append(a, True{})
		}
		if f > 0 {
			a = append(a, False{})
		}
		if len(a) == 1 {
			return a[0]
		}
		return a
	}
	switch x := a[0].(type) {
	case Eq:
		eqs = append(eqs, x)
	case Ne:
		nes = append(nes, x)
	case Ge:
		ges = append(ges, x)
	case Gt:
		gts = append(gts, x)
	case Le:
		les = append(les, x)
	case Lt:
		lts = append(lts, x)
	case As:
		ass = append(ass, x)
	case Na:
		nas = append(nas, x)
	case In:
		ins = append(ins, x)
	case Ni:
		nis = append(nis, x)
	case True:
		t++
	case False:
		f++
	default:
		panic(x)
	}
	return Unite(a[1:], eqs, nes, ges, gts, les, lts, ass, nas, ins, nis, t, f)
}

func Collapse(a And) Filter {
	return Unite(a, []Eq{}, []Ne{}, []Ge{}, []Gt{}, []Le{}, []Lt{}, []As{}, []Na{}, []In{}, []Ni{}, 0, 0)
}

func IsEmpty(f Filter) bool {
	switch v := f.(type) {
	case nil:
		return true

	case And:
		for _, x := range v {
			if !IsEmpty(x) {
				return false
			}
		}
		return true
	case Or:
		for _, x := range v {
			if !IsEmpty(x) {
				return false
			}
		}
		return true

	case In:
		if len(v) == 0 {
			return true
		}
		for _, x := range v {
			if len(x) > 0 {
				return false
			}
		}
		return true
	case Ni:
		if len(v) == 0 {
			return true
		}
		for _, x := range v {
			if len(x) > 0 {
				return false
			}
		}
		return true

	case Eq:
		return len(v) == 0
	case Ne:
		return len(v) == 0
	case Ge:
		return len(v) == 0
	case Gt:
		return len(v) == 0
	case Le:
		return len(v) == 0
	case Lt:
		return len(v) == 0
	case As:
		return len(v) == 0
	case Na:
		return len(v) == 0

	case True, False:
		return false

	case Formatter:
		return v.Size() == 0 && v.Len() == 0

	default:
		panic(v)
	}
}
