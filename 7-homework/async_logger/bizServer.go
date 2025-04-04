package main

import "context"

type bizServer struct {
	UnimplementedBizServer
}

func newBizServer() *bizServer {
	return &bizServer{}
}

func (b *bizServer) Check(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (b *bizServer) Add(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (b *bizServer) Test(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}
