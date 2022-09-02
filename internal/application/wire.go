package application

import (
	"github.com/google/wire"
	v1 "github.com/pinguo-icc/April/internal/application/v1"
	"github.com/pinguo-icc/April/internal/infrastructure/server"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(v1.OperationalPos), "*"),
	wire.Struct(new(v1.DataEnv), "*"),
	wire.Struct(new(v1.JsonConfig), "*"),
	wire.Struct(new(RouterDefines), "*"),
	// NewApp,
	wire.Bind(new(server.Register), new(*RouterDefines)),
)
