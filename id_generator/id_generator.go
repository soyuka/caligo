package id_generator

import (
	"github.com/matoous/go-nanoid"
)

func GetId() (string, error) {
	return gonanoid.Nanoid(14)
}
