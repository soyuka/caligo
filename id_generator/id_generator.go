package id_generator

import (
	"github.com/matoous/go-nanoid"
)

func GetId(idAlphabet string, idLength int) (string, error) {
	return gonanoid.Generate(idAlphabet, idLength)
}
