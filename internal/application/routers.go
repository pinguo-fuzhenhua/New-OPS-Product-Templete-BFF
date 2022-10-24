package application

import (
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	v1 "github.com/pinguo-icc/template/internal/application/v1"
	"github.com/pinguo-icc/template/internal/infrastructure/server"
)

type Context = khttp.Context

func PathParam(ctx Context, name string) (val string, ok bool) {
	raws := mux.Vars(ctx.Request())
	val, ok = raws[name]
	return
}

type RouterDefines struct {
	OPos    *v1.OperationalPos
	DataEnv *v1.DataEnv
	OpBasic *v1.JsonConfig
	Mpos    *v1.MaterialPositions
}

func (rd *RouterDefines) RouteRegister(r *khttp.Router) {
	var H = server.HandlerFunc

	v1 := r.Group("/v1")
	{
		v1.GET("/operational-positions", H(rd.OPos.PullByCodes))
		v1.GET("/json-config-show", H(rd.OpBasic.Show))
	}

	{
		// 数据环境
		v1.GET("/env", H(rd.DataEnv.ListEnv))
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/material-positions/{position}/categories", H(rd.Mpos.Categories))
		v2.GET("/material-positions/{position}/materials", H(rd.Mpos.Materials))
		v2.GET("/material-positions/material/detail", H(rd.Mpos.MaterialDetail))
	}
}
