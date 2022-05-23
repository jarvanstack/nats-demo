package main

import (
	nats "github.com/nats-io/nats.go"
)

var nc *nats.Conn

func init() {
	// Connect to a server
	nc, _ = nats.Connect(nats.DefaultURL)
}
