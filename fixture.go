// Package fixture provides all the functions to create a custom fixture
package fixture

import (
	"fmt"

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

	// add a data packet handler
	f.Receiver.AddDataHandler(1, f.BridgePacket)

	fmt.Println("Receiver type:", f.Receiver.ReceiverType())

	// configure the bridge
	// f.Bridge.ConfigFromFile("config.json")

	// start the receiver
	return f.Receiver.Listen()
}

func (f *Fixture) Stop() error {
	return f.Receiver.Stop()
}

func (f *Fixture) BridgePacket(data []byte) {
	f.Bridge.SendData(data)
}
