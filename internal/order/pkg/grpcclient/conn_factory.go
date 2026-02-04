package grpcclient

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

type ConnFactory struct {
	dis   *etcd.Registry
	conns map[string]*grpc.ClientConn
	mu    sync.Mutex
}

func NewConnFactory(dis *etcd.Registry) *ConnFactory {
	return &ConnFactory{
		dis:   dis,
		conns: make(map[string]*grpc.ClientConn),
	}
}

func (f *ConnFactory) Get(service string) *grpc.ClientConn {
	f.mu.Lock()
	defer f.mu.Unlock()

	if conn, ok := f.conns[service]; ok {
		return conn
	}

	conn, err := kratosgrpc.DialInsecure(
		context.Background(),
		kratosgrpc.WithEndpoint("discovery:///"+service),
		kratosgrpc.WithDiscovery(f.dis),
	)
	if err != nil {
		panic(err)
	}

	f.conns[service] = conn
	return conn
}
