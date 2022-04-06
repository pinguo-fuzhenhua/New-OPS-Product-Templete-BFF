package application

import (
	"github.com/gorilla/mux"
	v1 "github.com/pinguo-icc/Camera360/internal/application/v1"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/server"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type Context = khttp.Context

func PathParam(ctx Context, name string) (val string, ok bool) {
	raws := mux.Vars(ctx.Request())
	val, ok = raws[name]
	return
}

type RouterDefines struct {
	OPos *v1.OperationalPos
	DataEnv *v1.DataEnv
}

func (rd *RouterDefines) RouteRegister(r *khttp.Router) {
	var H = server.HandlerFunc

	v1 := r.Group("/v1")
	{
		v1.GET("/operational-positions", H(rd.OPos.PullByCodes))
	}

	{
		// 数据环境
		v1.GET("/env", H(rd.DataEnv.ListEnv))
	}
}
