package receiver

import (
	"fmt"
	"log"
	"net"

	"gitlab.com/patopest/go-sacn"
	"gitlab.com/patopest/go-sacn/packet"
)

type SacnReceiver struct {
	receiverType string
	receiver     *sacn.Receiver
	universes    []*Universe
	stopChan     chan bool
	activeUnis   map[uint16]bool
}

// Universe is a struct that represents a sACN Universe. It listens to a set of channels of a sACN Universe.
type Universe struct {
	universe  uint16            // the universe number
	nChannels uint16            // the number of channels in the universe (aka the size of the universe)
	channel   uint16            // the start channel, the fixture will ignore everything before this channel
	fChannels uint16            // the number of channels that the fixture will listen to, the fixture will ignore everything after (channel + fChannels - 1)
	callback  func(data []byte) // the callback that will be called when a packet is received
	connected bool              // a flag that indicates if a connection has already been established (default is false)
}

func newSacnReceiver() *SacnReceiver {
	return &SacnReceiver{
		receiverType: "sacn",
	}
}

// NewUniverse creates a new universe structure that indicates the fixture should listen to a set of channels of a sACN Universe.
// To create a universe, you need to provide the universe number, the number of channels in the universe, the start channel, the number of channels that the fixture will listen to, and a callback function.
// The callback function will be called when a packet is received and the data will be passed as an argument.
// NewUniverse returns a pointer to the created Universe and any error encountered.
func newUniverse(universe, nChannels, channel, fChannels uint16, callback func(data []byte)) (*Universe, error) {

	// Check if the universe number is valid (between 1 and 63999)
	if universe < 1 || universe > 63999 {
		return nil, fmt.Errorf("invalid universe number: %d, valid range is 1-63999", universe)
	}

	// because sACN is very versatile, we can have a universe with more than 512 channels
	// however, because we are following the structure of a DMX fixture, we will limit the number of channels to 512
	// this may change in the future when we add fixture types that may support more than 512 channels
	if nChannels > 512 {
		return nil, fmt.Errorf("invalid number of channels: %d, valid range is 1-512 (this will be updated soon)", nChannels)
	}

	// Check if the channel is valid (between 1 and nChannels)
	if channel < 1 || channel > nChannels {
		return nil, fmt.Errorf("invalid channel number: %d; valid range is 1-%d", channel, nChannels)
	}

	// Check if the number of channels is valid (between 1 and nChannels)
	if fChannels < 1 || fChannels > nChannels {
		return nil, fmt.Errorf("invalid number of channels: %d; valid range is 1-%d", fChannels, nChannels)
	}

	// Check if the start channel and the number of channels are valid
	if (channel + fChannels - 1) > nChannels {
		return nil, fmt.Errorf("invalid start channel and number of channels: %d, %d; the sum should be less than or equal to %d", channel, fChannels, nChannels)
	}

	// Check if the callback function is nil
	if callback == nil {
		return nil, fmt.Errorf("callback function is nil")
	}

	return &Universe{
		universe:  universe,
		nChannels: nChannels,
		channel:   channel,
		fChannels: fChannels,
		callback:  callback,
		connected: false,
	}, nil
}

// AddUniverse adds a universe to the fixture.
// The function takes a pointer to a Universe object as an argument and returns any error encountered.
func (s *SacnReceiver) AddUniverse(u *Universe) error {
	// Check if the universe is nil
	if u == nil {
		return fmt.Errorf("universe is nil")
	}

	// Check if the universe is already in the list of universes
	if s.getUniverse(u.universe) != nil {
		// Not sure if this should be an error or just a warning or if we should just ignore it
		// return fmt.Errorf("universe %d is already in the list of universes", u.universe)

		// For now, let's just ignore it
		return nil
	}

	// Add the universe to the fixture
	s.universes = append(s.universes, u)

	return nil
}

func (s *SacnReceiver) ReceiverType() string {
	return s.receiverType
}

func (s *SacnReceiver) ConfigFromFile(filepath string) error {
	// Read the config from the file

	// for now just create a dummy config
	uni1, err := newUniverse(
		1, 512, 1, 7,
		func(data []byte) {
			log.Println("Data received:", data)
		},
	)
	if err != nil {
		log.Println(err)
	}

	return s.AddUniverse(uni1)
}

func (s *SacnReceiver) Listen() error {
	// Create a channel to stop the receiver
	s.stopChan = make(chan bool)

	// Create a map to keep track of the active universes
	s.activeUnis = make(map[uint16]bool)

	// Start listening for packets
	go s.loop()

	return nil
}

func (s *SacnReceiver) Stop() error {
	// Stop listening for packets
	s.receiver.Stop()

	// Close the stop channel
	close(s.stopChan)

	return nil
}

func (s *SacnReceiver) loop() {

	// Create interface instance
	iface, err := net.InterfaceByName("wlp2s0")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new sACN receiver
	s.receiver = sacn.NewReceiver(
		iface,
	)

	// Join the discovery universe
	s.receiver.JoinUniverse(sacn.DISCOVERY_UNIVERSE)

	// Add the universes to the receiver
	for _, uni := range s.universes {
		s.receiver.JoinUniverse(uni.universe)
	}

	// Add callbacks to the receiver
	s.receiver.RegisterPacketCallback(
		packet.PacketTypeDiscovery,
		s.handleDiscoveryPacket,
	)
	s.receiver.RegisterPacketCallback(
		packet.PacketTypeData,
		s.handleDataPacket,
	)
	s.receiver.RegisterPacketCallback(
		packet.PacketTypeSync,
		s.handleSyncPacket,
	)
	s.receiver.RegisterTerminationCallback(
		s.handleTerminationPacket,
	)

	// Start the receiver
	s.receiver.Start()

	fmt.Println("Listening for packets...")

	// Wait for the stop signal
	for {
		select {
		case <-s.stopChan:
			return
		default:
			continue
		}

	}
}

func (s *SacnReceiver) handleDiscoveryPacket(p packet.SACNPacket, source string) {
	d, ok := p.(*packet.DiscoveryPacket)
	if !ok {
		return
	}

	var i int
	for i = 0; i < d.GetNumUniverses(); i++ {

		// log
		log.Println("Discovered universe:", d.Universes[i])

		// Add the discovered universes to the active universes
		s.activeUnis[d.Universes[i]] = true

		// Get universe
		u := s.getUniverse(d.Universes[i])
		if u == nil {
			// Universe not found, that means this fixture has no interest in this universe, so let's continue
			continue
		}

		// Check if the fixture is already connected to the universe
		if !u.connected {
			// Join the universe
			s.receiver.JoinUniverse(u.universe)
			u.connected = true
		}
	}
}

func (s *SacnReceiver) handleDataPacket(p packet.SACNPacket, source string) {
	d, ok := p.(*packet.DataPacket)
	if !ok {
		return
	}

	// Get Data
	data := d.GetData()

	// Get Universe
	universe := s.getUniverse(d.Universe)

	// Call the callback function
	universe.callback(data[universe.channel:(universe.channel + universe.fChannels)])
}

func (s *SacnReceiver) handleSyncPacket(p packet.SACNPacket, source string) {
	// handle the sync packet
}

func (s *SacnReceiver) handleTerminationPacket(universe uint16) {
	// handle the universe packet
}

// Given a number, this function returns the universe that has that number or nil if the universe is not in the list of universes that the fixture is listening to.
func (s *SacnReceiver) getUniverse(universeNumber uint16) *Universe {
	// Check if the universe is in the list of universes that the fixture is listening to
	for _, u := range s.universes {
		if universeNumber == u.universe {
			return u
		}
	}
	return nil
}

// AddDataHandler adds a handler function that will be called when a packet is received.
func (s *SacnReceiver) AddDataHandler(universe uint16, handler func(data []byte)) {
	// Get the universe
	u := s.getUniverse(universe)

	// Check if the universe is nil
	if u == nil {
		log.Println("Universe not found")
		return
	}

	// Add the handler to the universe
	u.callback = handler
}
