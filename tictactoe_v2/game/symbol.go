package game

type Symbol string

const (
	X Symbol = "X"
	O Symbol = "O"
)

func (s Symbol) Opposite() Symbol {
	if s == X {
		return O
	}
	return X
}
