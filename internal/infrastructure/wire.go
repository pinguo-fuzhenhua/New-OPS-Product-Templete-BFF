package infrastructure

import (
	"github.com/google/wire"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/clientset"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/conf"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/server"
	ptrace "github.com/pinguo-icc/kratos-library/v2/trace"
)

var ProviderSet = wire.NewSet(
	conf.ProviderSet,
	ptrace.NewTracerProvider,
	clientset.NewClientSet,
	server.NewHttpServer,

	wire.FieldsOf(new(*clientset.ClientSet), "FieldDefinitionsClient"),
)
