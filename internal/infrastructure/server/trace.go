package server

import (
	"net"
	"net/http"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func traceFilter(provider trace.TracerProvider) khttp.FilterFunc {
	tracer := provider.Tracer("bff")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ping" {
				next.ServeHTTP(w, r)
				return
			}

			var pathTpl string
			if route := mux.CurrentRoute(r); route != nil {
				pathTpl, _ = route.GetPathTemplate()
			}
			if pathTpl == "" {
				pathTpl = r.URL.Path
			}

			ctx, span := tracer.Start(r.Context(), pathTpl, trace.WithSpanKind(trace.SpanKindServer))
			defer span.End()

			attr := semconv.HTTPServerAttributesFromHTTPRequest("April", pathTpl, r)
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil {
				if host != "" {
					attr = append(attr, semconv.NetPeerIPKey.String(host))
				}
			}

			span.SetAttributes(attr...)

			if tid := span.SpanContext().TraceID(); tid.IsValid() {
				w.Header().Set("X-Trace-Id", tid.String())
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
