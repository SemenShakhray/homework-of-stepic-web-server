package main

type bizServer struct {
	UnimplementedBizServer
}

func newBizServer() *bizServer {
	return &bizServer{}
}

func (b *bizServer) Check(*Nothing) *Nothing {
	return &Nothing{}
}

func (b *bizServer) Add(*Nothing) *Nothing {
	return &Nothing{}
}

func (b *bizServer) Test(*Nothing) *Nothing {
	return &Nothing{}
}
