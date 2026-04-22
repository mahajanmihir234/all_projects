package app

type Value string

const (
	X Value = "X"
	O Value = "O"
)

func (v Value) Opposite() Value {
	if v == X {
		return O
	}
	return X
}
