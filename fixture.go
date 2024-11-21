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

// type GenPacket struct {
// 	packetType string
// 	metadata   map[string]any
// 	data       []byte
// }

func LoadFromFile(path string) *Fixture {
	return &Fixture{
		Bridge:   bridge.NewBridge("grpc"),
		Receiver: receiver.NewReceiver("sacn"),
	}

}

func (f *Fixture) Start() error {
	// configure the receiver
	err := f.Receiver.ConfigFromFile("config.json")
	if err != nil {
		return err
	}

	// configure the bridge
	// f.Bridge.ConfigFromFile("config.json")

	// start the receiver
	return f.Receiver.Listen()
}
