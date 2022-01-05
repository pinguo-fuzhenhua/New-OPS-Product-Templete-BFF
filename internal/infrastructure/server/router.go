package server

import "github.com/go-kratos/kratos/v2/transport/http"

func registerRouter(e *http.Router) {
	// 自定义框架中 /ping 请求的处理函数
	// curl -I -X GET localhost:8000/ping
	e.GET("/ping", func(ctx http.Context) error {
		ctx.Response().WriteHeader(204)
		return nil
	})

	// 自定义 http 路由
}
