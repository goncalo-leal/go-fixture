// Package fixture provides all the functions to create a custom fixture
package fixture

import (
	"github.com/goncalo-leal/go-fixture/bridge"
	"github.com/goncalo-leal/go-fixture/receiver"
)

type Fixture struct {
	Bridge   bridge.Bridge
	Receiver receiver.Receiver
}

type GenPacket struct {
	packetType string
	metadata   map[string]any
	data       []byte
}

func loadFromFile(path string) {

}
