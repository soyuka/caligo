package handlers

import (
	c "github.com/soyuka/caligo/config"
	bolt "go.etcd.io/bbolt"
)

type Env struct {
	DB     *bolt.DB
	Config c.Config
}
