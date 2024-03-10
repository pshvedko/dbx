package builder

import "fmt"

type Comma int

func (c Comma) Format(f fmt.State, _ rune) {
	if c > 0 {
		_, _ = f.Write([]byte{','})
	}
}

type Holder int

func (h Holder) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprint(f, "$", int(h))
}

type Keyword string

func (k Keyword) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprint(f, string(k))
}

const (
	NULL    Keyword = "NULL"
	TRUE    Keyword = "TRUE"
	FALSE   Keyword = "FALSE"
	DEFAULT Keyword = "DEFAULT"
)
