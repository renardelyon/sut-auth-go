package application

import (
	"context"
	"sut-auth-go/lib/pkg/db"
	"sut-auth-go/lib/utils"

	"google.golang.org/grpc"
)

type Application struct {
	JwtWrapper  utils.JwtWrapper
	DbClients   db.Handler
	GrpcServer  *grpc.Server
	GrpcClients map[string]*grpc.ClientConn
	Context     context.Context
}
