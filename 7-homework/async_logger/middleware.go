package main

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type middleware struct {
	serverOptions []grpc.ServerOption
	auth          *authACL
	subs          *subscriber
}

func newMiddleware(auth *authACL, subs *subscriber) *middleware {
	mid := &middleware{
		auth: auth,
		subs: subs,
	}
	mid.serverOptions = []grpc.ServerOption{
		grpc.UnaryInterceptor(mid.unaryInterceptor),
		grpc.StreamInterceptor(mid.streamInterceptor),
	}

	return mid
}

func (m *middleware) Do(ctx context.Context, method string) error {
	var consumer, host string

	md, _ := metadata.FromIncomingContext(ctx)

	consumer = strings.Join(md.Get("consumer"), "")
	if p, ok := peer.FromContext(ctx); ok {
		host = p.Addr.String()
	}

	m.subs.Notify(&Event{
		Method:    method,
		Consumer:  consumer,
		Host:      host,
		Timestamp: time.Now().Unix(),
	})

	if !m.auth.Check(consumer, method) {
		return status.Errorf(codes.Unauthenticated, "failed authorization")
	}

	return nil
}

func (m *middleware) unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	err := m.Do(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func (m *middleware) streamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	err := m.Do(ss.Context(), info.FullMethod)
	if err != nil {
		return err
	}

	return handler(srv, ss)
}
