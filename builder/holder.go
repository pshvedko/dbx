package builder

import (
	"fmt"
	"strconv"
)

type FormatFunc func(fmt.State, rune)

func (f FormatFunc) Format(w fmt.State, r rune) {
	f(w, r)
}

const DefaultIntegerPageSize = 512

var integers [][]byte

func init() {
	b := make([]byte, 0, DefaultIntegerPageSize)
	var i, n, x int
	for cap(b) > x {
		b = strconv.AppendInt(b, int64(i), 10)
		integers = append(integers, b)
		b = b[len(b):]
		i++
		n, x = 10, 1
		for n < i {
			n *= 10
			x++
		}
	}
}
