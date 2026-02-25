package filter

type Set map[string]any
type Index map[int]any

func (s Set) Index(f Fielder) Index {
	index := make(Index, len(s))
	for key, value := range s {
		index[f.Index(key)] = value
	}
	return index
}

type Mutator interface {
	Projector
	WithPK(...string) Mutator
	WithSet(Set) Mutator
	WithValue(string, any) Mutator
}

type Wrapper struct {
	Projector
}

type WrapperWithValue struct {
	Projector
	index Index
}

func (w WrapperWithValue) Value(i int) (any, bool, bool) {
	v, ok := w.index[i]
	if ok {
		return v, false, false
	}
	return w.Projector.Value(i)
}

func (w Wrapper) WithSet(set Set) Mutator {
	return Wrapper{Projector: WrapperWithValue{Projector: w, index: set.Index(w)}}
}

func (w Wrapper) WithValue(key string, value any) Mutator {
	return Wrapper{Projector: WrapperWithValue{Projector: w, index: Index{w.Index(key): value}}}
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
