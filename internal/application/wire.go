package application

import (
	"github.com/google/wire"
	v1 "github.com/pinguo-icc/template/internal/application/v1"
	"github.com/pinguo-icc/template/internal/infrastructure/server"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(v1.OperationalPos), "*"),
	wire.Struct(new(v1.MaterialPositions), "*"),
	wire.Struct(new(v1.DataEnv), "*"),
	wire.Struct(new(v1.JsonConfig), "*"),
	wire.Struct(new(RouterDefines), "*"),
	// NewApp,
	wire.Bind(new(server.Register), new(*RouterDefines)),
)
