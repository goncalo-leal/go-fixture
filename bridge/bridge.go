// Purpose: This file contains the bridge package which is responsible for bridging the receiver and the actuator.
package bridge

type Bridge interface {
	BridgeType() string                   // BridgeType returns the type of the bridge
	ConfigFromFile(filepath string) error // ConfigFromFile reads the config from the file
	Start() error                         // Start launches a thread
	Stop() error                          // Stop stops the thread
	SendData(data []byte) error           // SendData sends data to the actuator
}

func NewBridge(bridgeType string) Bridge {
	switch bridgeType {
	case "grpc":
		return newGRPCBridge()
	default:
		return nil
	}
}
