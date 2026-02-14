package filter

type Set map[string]any

func (s Set) Index(f Fielder) map[int]any {
	index := make(map[int]any, len(s))
	for key, value := range s {
		i := f.Columns()[key]
		index[i] = value
	}
	return index
}

type Mutator interface {
	Projector
	WithPK(...string) Mutator
	WithValue(Set) Mutator
}

type Wrapper struct {
	Projector
}

type WrapperWithValue struct {
	Projector
	index map[int]any
}

func (w WrapperWithValue) Value(i int) (any, bool, bool) {
	v, ok := w.index[i]
	if ok {
		return v, false, false
	}
	return w.Projector.Value(i)
}

func (w Wrapper) WithValue(set Set) Mutator {
	return Wrapper{Projector: WrapperWithValue{Projector: w, index: set.Index(w)}}
}

type WrapperWithPK struct {
	Projector
	pk []string
}

func (w WrapperWithPK) PK() PK {
	return w.pk
}

func (w Wrapper) WithPK(pk ...string) Mutator {
	return Wrapper{Projector: WrapperWithPK{Projector: w, pk: pk}}
}

func NewProjector(j Projector) Mutator {
	return Wrapper{Projector: j}
}
