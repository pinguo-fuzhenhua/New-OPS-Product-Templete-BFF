package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	kerr "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/conf"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/cparam"
	"github.com/pinguo-icc/operational-basic-svc/pkg/denv"
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
		cparam.LogValuer,
	)

	var opts = []khttp.ServerOption{
		khttp.Logger(logger),
		khttp.Address(config.Address),
		khttp.Timeout(config.Timeout),
		khttp.ErrorEncoder(buildErrorEncoder(logger)),
		khttp.Middleware(
			recovery.Recovery(recovery.WithLogger(loggerWithMethod)),
			serverLogging(loggerWithMethod),
		),
		khttp.Filter(
			traceFilter(tracerProvider),
			cparam.Filter(),
			denv.HTTPFilter,
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

func buildErrorEncoder(logger log.Logger) khttp.EncodeErrorFunc {
	var errCodec = encoding.GetCodec(json.Name)

	return func(w http.ResponseWriter, r *http.Request, err error) {
		type errResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		se := kerr.FromError(err)
		resp := &errResp{Code: int(se.Code), Message: se.Message}
		body, merr := errCodec.Marshal(resp)
		if merr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithContext(r.Context(), logger).Log(log.LevelError, "b_err", err, "e_err", merr)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(int(se.Code))
		wlen, werr := w.Write(body)
		log.WithContext(r.Context(), logger).Log(log.LevelWarn, "b_err", err, "w_len", wlen, "w_err", werr)
	}
}
