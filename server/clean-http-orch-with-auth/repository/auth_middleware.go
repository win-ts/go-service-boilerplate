package repository

import (
	"context"

	"google.golang.org/grpc"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/protobuf/authPb"
)

type authMiddlewareRepository struct {
	authClient authPb.AuthGrpcServiceClient
}

// AuthMiddlewareRepositoryDependencies represents the dependencies for auth middleware repository
type AuthMiddlewareRepositoryDependencies struct {
	GRPCClient *grpc.ClientConn
}

// NewAuthMiddlewareRepository creates a new auth middleware repository
func NewAuthMiddlewareRepository(d AuthMiddlewareRepositoryDependencies) AuthMiddlewareRepository {
	authClient := authPb.NewAuthGrpcServiceClient(d.GRPCClient)
	return &authMiddlewareRepository{authClient}
}

func (r *authMiddlewareRepository) VerifyToken(ctx context.Context, token string) error {
	result, err := r.authClient.VerifyToken(ctx, &authPb.VerifyTokenReq{
		Token: token,
	})
	if err != nil {
		return err
	}
	if !result.Success {
		return err
	}

	return nil
}
