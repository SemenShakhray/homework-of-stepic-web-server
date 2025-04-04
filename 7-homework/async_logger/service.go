package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func StartMyMicroservice(ctx context.Context, addr, ACLData string) error {
	authACL, err := newAuth(ACLData)
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	subs := newSubscriber()
	mid := newMiddleware(authACL, subs)
	server := grpc.NewServer(mid.serverOptions...)

	RegisterBizServer(server, newBizServer())
	RegisterAdminServer(server, newAdminServer(subs))

	go server.Serve(listener)

	go func() {
		<-ctx.Done()
		subs.DettachAll()
		server.GracefulStop()
	}()

	return nil
}
