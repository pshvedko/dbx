package filter

type True struct{}

func (f True) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f True) To(b Builder, _ Projector) error {
	_, err := b.Eq(nil, true)
	return err
}

type False struct{}

func (f False) MarshalJSON() ([]byte, error) { return MarshalJSON(f) }

func (f False) To(b Builder, _ Projector) error {
	_, err := b.Eq(nil, false)
	return err
}

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
