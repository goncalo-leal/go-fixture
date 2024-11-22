package bridge

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/goncalo-leal/go-fixture/proto/data"
)

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

func (g *gRPCBridge) ConfigFromFile(filepath string) error {
	return nil
}

func (g *gRPCBridge) Start() error {
	return nil
}

func (g *gRPCBridge) Stop() error {
	return nil
}

func (g *gRPCBridge) SendData(data []byte) error {

	// Create a client connection to the server
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	// Close the connection when done
	defer conn.Close()

	client := pb.NewDataServiceClient(conn)

	// Call the DataCallback method on the server
	_, err = client.DataCallback(context.Background(), &pb.DataReceived{Data: data})
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	fmt.Println("Data sent:", data)

	return nil
}
