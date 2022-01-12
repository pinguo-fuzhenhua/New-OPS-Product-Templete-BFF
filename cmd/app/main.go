package main

import (
	"flag"
	"fmt"
	"os"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/conf"
	fdpkg "github.com/pinguo-icc/field-definitions/pkg"
	"github.com/pinguo-icc/kratos-library/v2/pdebug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var env string // 运行环境：dev, qa, prod

func init() {
	flag.StringVar(&env, "env", "prod", "specify runtime environment: dev, qa, prod")
}

func newApp(cfg *conf.Bootstrap, logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Logger(logger),
		kratos.Server(hs),
		kratos.Name(cfg.App.Name),
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

	klog := newLogger(env)
	defer klog.Sync()

	logger := log.With(klog, "trace_id", tracing.TraceID())

	app, cleanup, err := initApp(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fdpkg.ConfigureUploader(fdpkg.NewUploader(cfg.Qiniu))

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func newLogger(env string) *kzap.Logger {
	var zlogger *zap.Logger
	var err error
	switch env {
	case "prod", "qa":
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		cfg.EncoderConfig.MessageKey = ""
		zlogger, err = cfg.Build(
			zap.WithCaller(false),
			zap.AddStacktrace(zapcore.FatalLevel),
			zap.AddCallerSkip(3),
		)
	default:
		zlogger, err = zap.NewDevelopment(
			zap.WithCaller(false),
			zap.AddStacktrace(zapcore.FatalLevel),
			zap.AddCallerSkip(3),
		)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return kzap.NewLogger(zlogger)
}
