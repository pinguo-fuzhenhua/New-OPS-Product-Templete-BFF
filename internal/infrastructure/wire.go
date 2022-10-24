package infrastructure

import (
	"github.com/google/wire"
	ptrace "github.com/pinguo-icc/kratos-library/v2/trace"
	"github.com/pinguo-icc/template/internal/infrastructure/clientset"
	"github.com/pinguo-icc/template/internal/infrastructure/conf"
	"github.com/pinguo-icc/template/internal/infrastructure/server"
)

var ProviderSet = wire.NewSet(
	conf.ProviderSet,
	ptrace.NewTracerProvider,
	clientset.NewClientSet,
	server.NewHttpServer,
	wire.FieldsOf(new(*clientset.ClientSet), "FieldDefinitionsClient", "OperationalBasicClient", "MaterialPositionsClient"),
)
