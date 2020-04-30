package id_generator

import (
	"github.com/matoous/go-nanoid"
)

func GetId(idLength int) (string, error) {
	return gonanoid.Nanoid(idLength)
}
