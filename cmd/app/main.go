package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/kratos-library/v2/pdebug"
	"github.com/pinguo-icc/template/internal/infrastructure/conf"
)

var env string // 运行环境：dev, qa, prod

func init() {
	flag.StringVar(&env, "env", "prod", "specify runtime environment: dev, qa, prod")
}

func newApp(logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Logger(logger),
		kratos.Server(hs),
		kratos.Name("BFF:template"),
	)
}

func main() {
	flag.Parse()

	cfg, err := conf.Load(env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if env == "dev" {
		pdebug.Enable(true)
	}

	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	app, cleanup, err := initApp(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
