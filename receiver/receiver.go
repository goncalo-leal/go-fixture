package receiver

type Receiver interface {
	ReceiverType()
	ConfigFromFile(filepath string)
	Listen() // launches a thread
	Stop()
}
