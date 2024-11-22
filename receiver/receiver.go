package receiver

type Receiver interface {
	ReceiverType() string // returns the type of the receiver
	ConfigFromFile(filepath string) error
	Listen() error // launches a thread
	Stop() error
	AddDataHandler(universe uint16, handler func(data []byte))
}

func NewReceiver(receiverType string) Receiver {
	switch receiverType {
	case "sacn":
		return newSacnReceiver()
	default:
		return nil
	}
}
