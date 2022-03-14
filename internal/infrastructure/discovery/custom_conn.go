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
	conns  []CustomGRPCConn
	opts   []kgrpc.ClientOption
	logger *log.Helper
	locker *sync.RWMutex
	count  int64
}

func NewCustomConn(logger *log.Helper) *CustomConn {
	return &CustomConn{
		conns:  make([]CustomGRPCConn, 0),
		opts:   make([]kgrpc.ClientOption, 0),
		logger: logger,
		locker: &sync.RWMutex{},
	}
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
	newEndpoints := []string{}
	conns := make([]CustomGRPCConn, 0)
	for _, ins := range instances {
		for _, endpoint := range ins.Endpoints {
			tmp, _ := url.Parse(endpoint)
			conn, err := s.connect(context.Background(), kgrpc.WithEndpoint(tmp.Host))
			conns = append(conns, conn)
			if err != nil {
				s.logger.Errorf("connect server error,host=%s,message=%s", tmp.Host, err.Error())
			} else {
				newEndpoints = append(newEndpoints, endpoint)
			}
		}
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	s.conns = conns
}

func (s *CustomConn) Connect(ctx context.Context, opts ...kgrpc.ClientOption) error {
	conn, err := s.connect(context.TODO(), opts...)
	if err != nil {
		return err
	}
	s.conns = append(s.conns, conn)
	return nil
}

func (s *CustomConn) connect(ctx context.Context, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	clientOpts := append(s.opts, opts...)
	conn, err := kgrpc.DialInsecure(context.TODO(), clientOpts...)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *CustomConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	l := s.locker.RLocker()
	l.Lock()
	defer func() {
		l.Unlock()
	}()

	return s.pickup().Invoke(ctx, method, args, reply, opts...)
}

func (s *CustomConn) pickup() grpc.ClientConnInterface {
	length := len(s.conns)
	conn := s.conns[0]
	if length > 1 {
		x := atomic.AddInt64(&s.count, 1)
		idx := x % int64(length)
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
