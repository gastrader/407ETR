package client

import (
	"github.com/gastrader/407ETR/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	if err != nil{
		return nil,err
	}
	return &GRPCClient{
		Endpoint: endpoint,
		AggregatorClient: c,
	}, nil
}