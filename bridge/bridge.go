// Purpose: This file contains the bridge package which is responsible for bridging the receiver and the actuator.
package bridge

type Bridge interface {
	// BridgeType returns the type of the bridge
	BridgeType() string
}

func NewBridge(bridgeType string) Bridge {
	switch bridgeType {
	case "grpc":
		return newGRPCBridge()
	default:
		return nil
	}
}
