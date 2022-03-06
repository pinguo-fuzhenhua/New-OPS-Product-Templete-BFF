package discovery

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type CustomGRPCConn interface {
	io.Closer
	grpc.ClientConnInterface
}

type CustomConn struct {
	conns       []CustomGRPCConn
	serviceName string
	factory     func(opts ...kgrpc.ClientOption) (*grpc.ClientConn, error)
	endpoints   []string
}

func NewCustomConn() *CustomConn {
	return &CustomConn{
		conns: make([]CustomGRPCConn, 0),
	}
}

func (s *CustomConn) SetFactory(fn func(opts ...kgrpc.ClientOption) (*grpc.ClientConn, error)) {
	s.factory = fn
}
func (s *CustomConn) SetServiceName(n string) {
	s.serviceName = n
}
func (s *CustomConn) close() error {
	for _, c := range s.conns {
		c.Close()
	}
	return nil
}

func (s *CustomConn) Notify(instances []*registry.ServiceInstance) {
	go func() {
		time.Sleep(time.Second * 5)
		newEndpoints := []string{}
		for _, ins := range instances {
			for _, endpoint := range ins.Endpoints {
				fmt.Println(endpoint)
				tmp, _ := url.Parse(endpoint)
				fmt.Println(tmp.Host)
				err := s.Connect(kgrpc.WithEndpoint(tmp.Host))
				if err != nil {
					fmt.Println(err)
				} else {
					newEndpoints = append(newEndpoints, endpoint)
				}
			}
		}
		s.endpoints = newEndpoints
	}()
}

// Connect
// @TODO 初次启动的时候会连两次
func (s *CustomConn) Connect(opts ...kgrpc.ClientOption) error {
	conn, err := s.factory(opts...)
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		if conn.GetState() == connectivity.Ready {
			break
		}
		fmt.Println(conn.GetState().String())
		time.Sleep(time.Second)
	}
	s.conns = append(s.conns, conn)
	return nil
}

func (s *CustomConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	maxIdx := len(s.conns) - 1
	fmt.Println(s.conns[maxIdx].(*grpc.ClientConn).GetState())
	conn := s.conns[maxIdx]
	if ac := len(s.endpoints); ac > 1 {
		i := int(rand.Float64()*100) % ac
		idx := maxIdx - ac + 1 + i
		if idx > maxIdx {
			idx = maxIdx
		}
		conn = s.conns[idx]
	}
	return conn.Invoke(ctx, method, args, reply, opts...)
}

func (s CustomConn) Close() error {
	return s.close()
}

// NewStream begins a streaming RPC.
func (s *CustomConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return s.conns[0].NewStream(ctx, desc, method, opts...)
}
