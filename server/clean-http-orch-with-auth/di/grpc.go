package di

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	client *grpc.ClientConn
}

type grpcClientOptions struct {
	url string
}

func newGRPCClient(opts grpcClientOptions) (*grpcClient, error) {
	client, err := grpc.NewClient(
		opts.url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{client}, nil
}
