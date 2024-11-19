package bridge

type gRPCBridge struct {
	bridgeType string
}

func newGRPCBridge() *gRPCBridge {
	return &gRPCBridge{
		bridgeType: "grpc",
	}
}

func (g *gRPCBridge) BridgeType() string {
	return g.bridgeType
}
