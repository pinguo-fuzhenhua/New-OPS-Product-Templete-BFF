package server

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/conf"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/cparam"
	"go.opentelemetry.io/otel/trace"
)

type Register interface {
	RouteRegister(*khttp.Router)
}

// New new a bm server.
func NewHttpServer(config *conf.HTTP, tracerProvider trace.TracerProvider, logger log.Logger, r Register) (*khttp.Server, func()) {
	loggerWithMethod := log.With(
		logger,
		"method",
		log.Valuer(func(ctx context.Context) interface{} {
			if c, ok := ctx.(khttp.Context); ok {
				return c.Request().Method
			}
			return ""
		}),
		"path",
		log.Valuer(func(ctx context.Context) interface{} {
			if c, ok := ctx.(khttp.Context); ok {
				return c.Request().URL.Path
			}
			return ""
		}),
		"c_params",
		log.Valuer(func(ctx context.Context) interface{} {
			if p := cparam.FromContext(ctx); p != nil {
				return fmt.Sprintf("%+v", p)
			}
			return ""
		}),
	)

	var opts = []khttp.ServerOption{
		khttp.Logger(logger),
		khttp.Address(config.Address),
		khttp.Timeout(config.Timeout),
		khttp.Middleware(
			recovery.Recovery(recovery.WithLogger(loggerWithMethod)),
			logging.Server(loggerWithMethod),
		),
		khttp.Filter(
			traceFilter(tracerProvider),
			cparam.Filter(),
		),
	}

	svc := khttp.NewServer(opts...)
	route := svc.Route("/")
	registerRouter(route)
	r.RouteRegister(route)

	cancelFn := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		svc.Shutdown(ctx)
	}

	return svc, cancelFn
}
