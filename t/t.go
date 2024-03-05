package help

type Object struct {
	Id      uint32  `json:"id"`
	Bool    bool    `json:"o_bool,omitempty"`
	Float32 float32 `json:"o_float_32,omitempty"`
	Float64 float64 `json:"o_float_64,omitempty"`
	Int     int     `json:"o_int,omitempty"`
	Int16   int16   `json:"o_int_16,omitempty"`
	Null    any     `json:"o_null,omitempty"`
	String  string  `json:"o_string,omitempty"`
}

func (o Object) Table() string {
	return "objects"
}

func (o Object) Names() []string {
	return []string{"id", "o_bool", "o_float_32", "o_float_64", "o_int", "o_int_16", "o_null", "o_string"}
}

func (o *Object) Values() []any {
	return []any{&o.Id, &o.Bool, &o.Float32, &o.Float64, &o.Int, &o.Int16, &o.Null, &o.String}
}
