package filter

type Tree struct {
	Right, Left *Tree
	Value       int
}

func (t *Tree) Height() int {
	if t == nil {
		return 0
	}
	r := t.Right.Height()
	l := t.Left.Height()
	if r > l {
		return r + 1
	}
	return l + 1
}
