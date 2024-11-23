package repository

import (
	"context"

	"google.golang.org/grpc"

	authPb "github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/protobuf/auth"
)

type authMiddlewareRepository struct {
	authClient authPb.AuthGrpcServiceClient
}

// AuthMiddlewareRepositoryConfig represents the configuration for auth middleware repository
type AuthMiddlewareRepositoryConfig struct {
}

// AuthMiddlewareRepositoryDependencies represents the dependencies for auth middleware repository
type AuthMiddlewareRepositoryDependencies struct {
	GrpcClient *grpc.ClientConn
}

// NewAuthMiddlewareRepository creates a new auth middleware repository
func NewAuthMiddlewareRepository(c AuthMiddlewareRepositoryConfig, d AuthMiddlewareRepositoryDependencies) AuthMiddlewareRepository {
	_ = c
	authClient := authPb.NewAuthGrpcServiceClient(d.GrpcClient)
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
