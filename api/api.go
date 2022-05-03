package api

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/oraclemgr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	npool.UnimplementedOracleManagerServer
}

func Register(server grpc.ServiceRegistrar) {
	npool.RegisterOracleManagerServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return npool.RegisterOracleManagerHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
