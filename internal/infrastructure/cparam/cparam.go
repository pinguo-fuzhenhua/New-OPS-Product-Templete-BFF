package cparam

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	klog "github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	baseparam "github.com/pinguo-icc/go-base/v2/param"
)

// 公共参数

type Params = baseparam.CommonParam

func New(r *http.Request) *Params {
	p := make(map[string]string, 32)
	for k, v := range r.Header {
		if strings.HasPrefix(k, "Pg-") || strings.HasPrefix(k, "PG-") || strings.HasPrefix(k, "pg-") {
			if len(v) > 0 {
				p[k] = v[0]
			}
		}
	}

	return baseparam.New(p)
}

type ctxKey struct{}

func Store(ctx context.Context, p *Params) context.Context {
	return context.WithValue(ctx, ctxKey{}, p)
}

func FromContext(ctx context.Context) *Params {
	v := ctx.Value(ctxKey{})
	if v != nil {
		return v.(*Params)
	}
	return nil
}

// Filter create Params and store
func Filter() khttp.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			r = r.WithContext(Store(r.Context(), New(r)))
			next.ServeHTTP(rw, r)
		})
	}
}

var LogValuer = logValuer()

func logValuer() klog.Valuer {
	return func(ctx context.Context) interface{} {
		if p := FromContext(ctx); p != nil {
			return fmt.Sprintf("%+v", p)
		}
		return ""
	}
}
