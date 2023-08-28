package client

import (
	"context"
	"fmt"

	"github.com/gastrader/407ETR/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	fmt.Println("The endpoint is:", endpoint)
	conn, err := grpc.Dial("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	if err != nil{
		return nil,err
	}
	return &GRPCClient{
		Endpoint: endpoint,
		client: c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}