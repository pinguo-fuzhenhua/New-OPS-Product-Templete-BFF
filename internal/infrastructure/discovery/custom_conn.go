package discovery

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

type CustomGRPCConn interface {
	io.Closer
	grpc.ClientConnInterface
}

type CustomConn struct {
	conns  []CustomGRPCConn
	opts   []kgrpc.ClientOption
	logger *log.Helper
	notify chan []*registry.ServiceInstance
	conn   chan CustomGRPCConn
}

func NewCustomConn(logger *log.Helper) *CustomConn {
	x := &CustomConn{
		opts:   make([]kgrpc.ClientOption, 0),
		logger: logger,
		notify: make(chan []*registry.ServiceInstance, 10),
		conn:   make(chan CustomGRPCConn, 2),
	}
	go x.watch()
	return x
}

func (s *CustomConn) WithKGRPCClientOption(opts ...kgrpc.ClientOption) {
	s.opts = opts
}

func (s *CustomConn) close() error {
	for _, c := range s.conns {
		c.Close()
	}
	return nil
}

func (s *CustomConn) Notify(instances []*registry.ServiceInstance) {
	s.notify <- instances
}

func (s *CustomConn) pickup(ctx context.Context, conns []CustomGRPCConn) {
	for {
		if len(conns) == 0 {
			time.Sleep(time.Second)
		}
		for _, conn := range conns {
			select {
			case <-ctx.Done():
				return
			default:
				s.conn <- conn
			}
		}
	}
}
func (s *CustomConn) watch() {
	ctx, cancel := context.WithCancel(context.Background())
	go s.pickup(ctx, s.conns)
	update := func(instances []*registry.ServiceInstance) {
		conns := make([]CustomGRPCConn, 0)
		for _, ins := range instances {
			for _, endpoint := range ins.Endpoints {
				tmp, _ := url.Parse(endpoint)
				conn, err := s.connect(context.Background(), kgrpc.WithEndpoint(tmp.Host))
				if err != nil {
					s.logger.Errorf("connect server error,host=%s,message=%s", tmp.Host, err.Error())
				} else {
					conns = append(conns, conn)
				}
			}
		}
		cancel()
		ctx, cancel = context.WithCancel(context.Background())
		total := len(s.conn)
		go s.pickup(ctx, conns)
		go func() {
			t := time.After(time.Second * 30)
			for i := 0; i < total; i++ {
				select {
				case <-s.conn:
				case <-t:
				}
			}
		}()
		s.conns = conns
	}
	for instances := range s.notify {
		update(instances)
	}
}

func (s *CustomConn) Connect(ctx context.Context, opts ...kgrpc.ClientOption) error {
	conn, err := s.connect(ctx, opts...)
	if err != nil {
		return err
	}
	s.conns = append(s.conns, conn)
	return nil
}

func (s *CustomConn) connect(ctx context.Context, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	clientOpts := append(s.opts, opts...)
	conn, err := kgrpc.DialInsecure(ctx, clientOpts...)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *CustomConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	c := <-s.conn
	return c.Invoke(ctx, method, args, reply, opts...)
}

func (s *CustomConn) Close() error {
	return s.close()
}

// NewStream begins a streaming RPC.
func (s *CustomConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	c := <-s.conn
	return c.NewStream(ctx, desc, method, opts...)
}
