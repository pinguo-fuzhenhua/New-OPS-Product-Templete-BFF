package server

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/pinguo-icc/template/internal/infrastructure/conf"
)

type Register interface {
	RouteRegister(*khttp.Router)
}

// New new a bm server.
func NewHttpServer(config *conf.HTTP, logger log.Logger, r Register) (*khttp.Server, func()) {
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
			cors(),
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

func cors() khttp.FilterFunc {
	return handlers.CORS(
		handlers.AllowedOriginValidator(func(s string) bool {
			return strings.Contains(s, "camera360.com")
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "PUT", "PATCH", "DELETE"}),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{
			"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Authorization",
		}),
	)
}
