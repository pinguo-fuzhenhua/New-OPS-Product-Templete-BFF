package discovery

import (
	"context"
	"io"
	"net/url"
	"sync"
	"sync/atomic"

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
	conns         []CustomGRPCConn
	opts          []kgrpc.ClientOption
	logger        *log.Helper
	endpointCount int
	offset        int
	count         int64
	notify        chan []*registry.ServiceInstance
	locker        *sync.RWMutex
}

func NewCustomConn(logger *log.Helper) *CustomConn {
	x := &CustomConn{
		conns:  make([]CustomGRPCConn, 0),
		opts:   make([]kgrpc.ClientOption, 0),
		logger: logger,
		notify: make(chan []*registry.ServiceInstance, 10),
		locker: &sync.RWMutex{},
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

func (s *CustomConn) watch() {
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

		if len(s.conns) > 50 {
			s.locker.Lock()
			defer s.locker.Unlock()
			s.conns = conns
			s.offset = 0
		} else {
			s.conns = append(s.conns, conns...)
			s.offset = len(s.conns)
		}
		s.endpointCount = len(conns)
		s.logger.Debugf("offset=%v, endpointcount=%v", s.offset, s.endpointCount)
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
	return s.pickup().Invoke(ctx, method, args, reply, opts...)
}

func (s *CustomConn) pickup() grpc.ClientConnInterface {
	s.locker.RLock()
	defer s.locker.RUnlock()

	conn := s.conns[s.offset]
	if s.endpointCount > 1 {
		x := atomic.AddInt64(&s.count, 1)
		idx := int(x)%s.endpointCount + s.offset
		conn = s.conns[idx]
	}
	return conn
}

func (s *CustomConn) Close() error {
	return s.close()
}

// NewStream begins a streaming RPC.
func (s *CustomConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return s.pickup().NewStream(ctx, desc, method, opts...)
}
