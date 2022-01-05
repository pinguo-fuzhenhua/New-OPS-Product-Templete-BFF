package server

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/render"
	"github.com/pinguo-icc/go-base/v2/iuser"
)

func HandlerFunc(handler interface{}) khttp.HandlerFunc {
	// 使日志中的 operation 字段为请求 path
	setOperation := func(ctx context.Context) {
		if tr, ok := transport.FromServerContext(ctx); ok {
			if tr, ok := tr.(khttp.Transporter); ok {
				khttp.SetOperation(ctx, tr.PathTemplate())
			}
		}
	}

	switch handle := handler.(type) {
	case func(khttp.Context) error:
		return func(ctx khttp.Context) error {
			setOperation(ctx)
			next := ctx.Middleware(func(mCtx context.Context, _ interface{}) (interface{}, error) {
				err := handle(ctx)
				return nil, err
			})

			_, err := next(ctx, nil)
			return err
		}
	case func(khttp.Context) (interface{}, error):
		return func(ctx khttp.Context) error {
			setOperation(ctx)
			next := ctx.Middleware(func(mCtx context.Context, _ interface{}) (interface{}, error) {
				return handle(ctx)
			})

			res, err := next(ctx, nil)
			if err != nil {
				return err
			}
			return render.RenderJSON(ctx, res)
		}
	case func(http.ResponseWriter, *http.Request):
		return func(ctx khttp.Context) error {
			setOperation(ctx)
			next := ctx.Middleware(func(mCtx context.Context, _ interface{}) (interface{}, error) {
				req := ctx.Request()
				if mCtx != ctx {
					req = req.WithContext(mCtx)
				}
				handle(ctx.Response(), req)
				return nil, nil
			})

			_, _ = next(ctx, nil)
			return nil
		}
	default:
		panic("invalid handler")
	}
}

// LoginRequired http filter for backend system
func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := iuser.AcquireUserFromCookie(r)
		if err != nil {
			render.LoginRequired(w)
			return
		}

		ctx := iuser.WithUser(r.Context(), user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
