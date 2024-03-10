package filter

func Deep[T interface {
	comparable
	Right() T
	Left() T
}](t T) int {
	var z T
	if t == z {
		return 0
	}
	r := Deep(t.Right())
	l := Deep(t.Left())
	if r > l {
		return r + 1
	}
	return l + 1
}

func Loop[I interface {
	comparable
	Next() I
}](i I) bool {
	var (
		z I
		h = i
		b bool
	)
	for i != z {
		i = i.Next()
		if h == i {
			return true
		}
		b = !b
		if b {
			h = h.Next()
		}
	}
	return false
}
