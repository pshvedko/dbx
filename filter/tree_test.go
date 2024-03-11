package filter_test

import (
	"testing"

	"github.com/pshvedko/dbx/filter"
)

type tree struct {
	value       any
	right, left *tree
}

func (t *tree) Left() *tree {
	return t.left
}

func (t *tree) Right() *tree {
	return t.right
}

func TestDeep(t *testing.T) {
	tests := []struct {
		name string
		tree tree
		want int
	}{
		// TODO: Add test cases.
		{
			//     1
			//    / \
			//   2   3
			//  / \
			// 5   4
			//    /
			//   6
			name: "",
			tree: tree{
				value: 1,
				right: &tree{
					value: 3,
					right: nil,
					left:  nil,
				},
				left: &tree{
					value: 2,
					right: &tree{
						value: 4,
						right: nil,
						left: &tree{
							value: 6,
							right: nil,
							left:  nil,
						},
					},
					left: &tree{
						value: 5,
						right: nil,
						left:  nil,
					},
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter.Deep(&tt.tree); got != tt.want {
				t.Errorf("Deep() = %v, want %v", got, tt.want)
			}
		})
	}
	var tr Tree = tree1{
		r: &tree1{
			r: &tree1{
				r: &tree1{},
				l: nil,
			},
			l: nil,
		},
		l: nil,
	}
	if got := filter.Deep(tr); got != 4 {
		t.Errorf("Deep() = %v, want %v", got, 4)
	}
}

type tree1 struct {
	r, l *tree1
}

func (t tree1) Left() Tree {
	if t.l != nil {
		return t.l
	}
	return nil
}

func (t tree1) Right() Tree {
	if t.r != nil {
		return t.r
	}
	return nil
}

type Tree interface {
	Left() Tree
	Right() Tree
}

type list struct {
	value any
	next  *list
}

func (l *list) Next() *list {
	return l.next
}

func TestLoop(t *testing.T) {
	tests := []struct {
		name string
		list list
		want bool
	}{
		// TODO: Add test cases.
		{
			// 1
			name: "",
			list: list{
				value: 1,
				next:  nil,
			},
			want: false,
		},
		{
			// 1-2-3-4-5
			name: "",
			list: list{
				value: 1,
				next: &list{
					value: 2,
					next: &list{
						value: 3,
						next: &list{
							value: 4,
							next: &list{
								value: 5,
								next:  nil,
							},
						},
					},
				},
			},
			want: false,
		},
		{
			// 1-1
			name: "",
			list: func() list {
				l := list{
					value: 1,
					next:  nil,
				}
				l.next = &l
				return l
			}(),
			want: true,
		},
		{
			// 1-2-1
			name: "",
			list: func() list {
				l := list{
					value: 1,
					next: &list{
						value: 2,
						next:  nil,
					},
				}
				l.next.next = &l
				return l
			}(),
			want: true,
		},
		{
			// 1-2-3-2
			name: "",
			list: func() list {
				l := list{
					value: 1,
					next: &list{
						value: 2,
						next: &list{
							value: 3,
							next:  nil,
						},
					},
				}
				l.next.next.next = l.next
				return l
			}(),
			want: true,
		},
		{
			// 1-2-3-...-10-...-100-10
			name: "",
			list: func() list {
				l := &list{value: 100}
				c := l
				i := 100
				for i > 10 {
					i--
					l = &list{value: i, next: l}
				}
				x := l
				for i > 1 {
					i--
					l = &list{value: i, next: l}
				}
				c.next = x
				return *l
			}(),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter.Loop(&tt.list); got != tt.want {
				t.Errorf("Loop() = %v, want %v", got, tt.want)
			}
		})
	}
	if !filter.Loop(list100(1)) {
		t.Errorf("Loop() = %v, want %v", false, true)
	}
}

type list100 int

func (l list100) Next() list100 {
	return l%100 + 1
}
