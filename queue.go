package queue

import (
	"errors"
	"net/http"
)

var (
	ErrNoConnection = errors.New("no connection")
	ErrNoConsumers  = errors.New("no consumers")
	ErrNoMessage    = errors.New("no message")
)

var defaultConnection *Conn

type Options struct {
	Name    string
	Address string
}

// Init initializes the default connection to the message queue server.
func Init(opts Options) {
	c := &Conn{
		httpClient: &http.Client{},
		address:    opts.Address,
		name:       opts.Name,
	}

	defaultConnection = c
}

// Connect creates a new connection to the message queue server.
func Connect(opts Options) *Conn {
	c := &Conn{
		httpClient: &http.Client{},
		address:    opts.Address,
	}

	return c
}

type Conn struct {
	httpClient *http.Client
	address    string
	name       string
}

func Ready() bool {
	if defaultConnection == nil {
		return false
	}

	return defaultConnection.Ready()
}

func (cn *Conn) Ready() bool {
	if cn == nil {
		return false
	}

	return true
}
