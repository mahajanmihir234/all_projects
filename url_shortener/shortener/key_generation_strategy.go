package shortener

import (
	"math/rand"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type KeyGenerationStrategy interface {
	GenerateKey(id int) string
}

type RandomGenerationStrategy struct {
	characters string
	keyLength  int
}

func (r RandomGenerationStrategy) GenerateKey(id int) string {
	key := ""
	for range r.keyLength {
		key += string(r.characters[rand.Intn(len(r.characters))])
	}
	return key
}

func NewRandomStrategy() RandomGenerationStrategy {
	return RandomGenerationStrategy{
		characters: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		keyLength:  6,
	}
}

type UUIDStrategy struct {
	keyLength int
}

func (s UUIDStrategy) GenerateKey(id int) string {
	uuidString := uuid.New().String()
	uuidString = strings.ReplaceAll(uuidString, "-", "")
	return uuidString[:s.keyLength]
}

func NewUUIDStrategy() UUIDStrategy {
	return UUIDStrategy{
		keyLength: 6,
	}
}

type Base62Strategy struct {
	base62Characters  string
	base              int
	min6ChartIdOffset int
}

func (s Base62Strategy) GenerateKey(id int) string {
	if id == 0 {
		return string(s.base62Characters[0])
	}

	idWithOffset := id + s.min6ChartIdOffset
	result := []string{}
	for idWithOffset > 0 {
		result = append(result, string(s.base62Characters[idWithOffset%s.base]))
		idWithOffset /= s.base
	}
	slices.Reverse(result)
	return strings.Join(result, "")
}

func NewBase62Strategy() Base62Strategy {
	return Base62Strategy{
		base62Characters:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		base:              62,
		min6ChartIdOffset: 62 ^ 5,
	}
}
