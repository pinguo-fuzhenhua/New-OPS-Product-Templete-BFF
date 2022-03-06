package discovery

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
)

type CustomGRPCConn interface {
	io.Closer
	grpc.ClientConnInterface
}

type CustomConn struct {
	conn        CustomGRPCConn
	conns       []CustomGRPCConn
	serviceName string
	factory     func() (CustomGRPCConn, error)
}

func NewCustomConn() *CustomConn {
	return &CustomConn{
		conns: make([]CustomGRPCConn, 0),
	}
}

func (s *CustomConn) SetFactory(fn func() (CustomGRPCConn, error)) {
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

func (s *CustomConn) Notify() {
	s.Connect(false)
}

// Connect
// @TODO 初次启动的时候会连两次
func (s *CustomConn) Connect(isInit bool) error {
	if !isInit && s.conn == nil {
		return nil
	}
	conn, err := s.factory()
	if err != nil {
		return err
	}
	s.conn = conn
	s.conns = append(s.conns, conn)
	if len(s.conns) > 1 {
		max := len(s.conns) - 1
		for i := 0; i < max; i++ {
			fmt.Println("close old connection")
			go func(cn CustomGRPCConn) {
				time.Sleep(time.Second * 60)
				cn.Close()
			}(s.conns[i])
		}
		s.conns = s.conns[max:]
	}
	return nil
}

func (s *CustomConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return s.conn.Invoke(ctx, method, args, reply, opts...)
}

func (s CustomConn) Close() error {
	return s.close()
}

// NewStream begins a streaming RPC.
func (s *CustomConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return s.conn.NewStream(ctx, desc, method, opts...)
}
