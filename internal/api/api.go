package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brocaar/loraserver/api"
	"github.com/brocaar/loraserver/internal/loraserver"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// GetGRPCServer returns the gRPC API handler.
func GetGRPCServer(ctx context.Context, lsCtx loraserver.Context) *grpc.Server {
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	pb.RegisterApplicationServer(server, NewApplicationAPI(lsCtx))
	pb.RegisterNodeServer(server, NewNodeAPI(lsCtx))
	pb.RegisterChannelListServer(server, NewChannelListAPI(lsCtx))
	pb.RegisterChannelServer(server, NewChannelAPI(lsCtx))
	return server
}

// GetJSONGateway returns the JSON gateway for the gRPC API.
func GetJSONGateway(ctx context.Context, lsCtx loraserver.Context, grpcBind string) (http.Handler, error) {
	bindParts := strings.SplitN(grpcBind, ":", 2)
	if len(bindParts) != 2 {
		return nil, errors.New("get port from http-bind failed")
	}
	apiEndpoint := fmt.Sprintf("localhost:%s", bindParts[1])

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	mux := runtime.NewServeMux()

	if err := pb.RegisterApplicationHandlerFromEndpoint(ctx, mux, apiEndpoint, opts); err != nil {
		return nil, fmt.Errorf("register application handler error: %s", err)
	}
	if err := pb.RegisterNodeHandlerFromEndpoint(ctx, mux, apiEndpoint, opts); err != nil {
		return nil, fmt.Errorf("register node handler error: %s", err)
	}
	if err := pb.RegisterChannelListHandlerFromEndpoint(ctx, mux, apiEndpoint, opts); err != nil {
		return nil, fmt.Errorf("register channel-list handler error: %s", err)
	}
	if err := pb.RegisterChannelHandlerFromEndpoint(ctx, mux, apiEndpoint, opts); err != nil {
		return nil, fmt.Errorf("register channel handler error: %s", err)
	}

	return mux, nil
}