package handlers

import (
	t "github.com/soyuka/caligo/transports"
	c "github.com/soyuka/caligo/config"
)


type Env struct {
	Transport t.Transport
	Config c.Config
}
