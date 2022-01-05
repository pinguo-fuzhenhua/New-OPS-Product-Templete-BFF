package clientset

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pinguo-icc/template/internal/infrastructure/conf"
)

// ClientSet gRPC Client Set
type ClientSet struct {
	// 注入其他微服务客户端
}

// NewClientSet new gRPC Client Set
func NewClientSet(c *conf.Params, logger log.Logger) (*ClientSet, func(), error) {
	conn, err := kgrpc.DialInsecure(
		context.TODO(),
		kgrpc.WithEndpoint(c.ArticleSvcAddr),
		kgrpc.WithMiddleware(
			logging.Client(logger),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	cancelFn := func() {
		conn.Close()
	}

	_ = conn

	h := &ClientSet{
		// 连接其他微服务客户端
	}

	return h, cancelFn, nil
}
