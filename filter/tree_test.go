package filter

import "testing"

func TestTree_Height(t *testing.T) {
	tests := []struct {
		name string
		root Tree
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
			root: Tree{
				Right: &Tree{
					Right: nil,
					Left:  nil,
					Value: 3,
				},
				Left: &Tree{
					Right: &Tree{
						Right: nil,
						Left: &Tree{
							Right: nil,
							Left:  nil,
							Value: 6,
						},
						Value: 4,
					},
					Left: &Tree{
						Right: nil,
						Left:  nil,
						Value: 5,
					},
					Value: 2,
				},
				Value: 1,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			if got := tt.root.Height(); got != tt.want {
				t1.Errorf("Height() = %v, want %v", got, tt.want)
			}
		})
	}
}
