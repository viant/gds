package redblack

// Integer implements Getter
type Integer int

// Get returns int value
func (m Integer) Get() int {
	return int(m)
}

// String  implements Getter
type String string

// Get returns string value
func (m String) Get() string {
	return string(m)
}

// Float64  implements Getter
type Float64 float64

// Get returns float64 value
func (m Float64) Get() float64 {
	return float64(m)
}

type Complex128 complex128

func (c Complex128) Get() complex128 {
	return complex128(c)
}
