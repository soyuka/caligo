package transports

import (
	"errors"
	"fmt"
	"net/url"
	c "github.com/soyuka/caligo/config"
)

var (
	// ErrInvalidTransportDSN is returned when the Transport's DSN is invalid.
	ErrInvalidTransportDSN = errors.New("invalid transport DSN")
	// ErrClosedTransport is returned by the Transport's Dispatch and AddSubscriber methods after a call to Close.
	ErrClosedTransport = errors.New("hub: read/write on closed Transport")
)

type Transport interface {
	Put(id string, url string) error
	Get(id string) (string, error)
	Count() (int64, error)
}

// NewTransport create a transport using the backend matching the given TransportURL.
func NewTransport(config *c.Config) (Transport, error) {
	u, err := url.Parse(config.DB)
	if err != nil {
		return nil, fmt.Errorf("transport_url: %w", err)
	}

	switch u.Scheme {
	case "bolt":
		return NewBoltTransport(u)
	case "redis":
		return NewRedisTransport(u)
	}

	return nil, fmt.Errorf("%q: no such transport available: %w", config.DB, ErrInvalidTransportDSN)
}
