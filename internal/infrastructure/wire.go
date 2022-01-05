package infrastructure

import (
	"github.com/google/wire"
	"github.com/pinguo-icc/template/internal/infrastructure/clientset"
	"github.com/pinguo-icc/template/internal/infrastructure/conf"
	"github.com/pinguo-icc/template/internal/infrastructure/server"
)

var ProviderSet = wire.NewSet(
	conf.ProviderSet,
	clientset.NewClientSet,
	server.NewHttpServer,
	// wire.FieldsOf(new(*clientset.ClientSet), "FooClient"),
)
