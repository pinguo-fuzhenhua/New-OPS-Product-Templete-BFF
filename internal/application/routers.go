package application

import (
	"github.com/gorilla/mux"
	"github.com/pinguo-icc/template/internal/infrastructure/server"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type Context = khttp.Context

func PathParam(ctx Context, name string) (val string, ok bool) {
	raws := mux.Vars(ctx.Request())
	val, ok = raws[name]
	return
}

type RouterDefines struct {
	E *Example
}

func (rd *RouterDefines) RouteRegister(r *khttp.Router) {
	var H = server.HandlerFunc

	v1 := r.Group("/v1")
	{
		{
			b := v1.Group("/example")

			b.GET("/", H(rd.E.Get))
			b.POST("/", H(rd.E.Post))
		}
	}
}
