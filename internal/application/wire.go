package application

import (
	"github.com/google/wire"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/server"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Example), "*"),
	wire.Struct(new(RouterDefines), "*"),
	// NewApp,
	wire.Bind(new(server.Register), new(*RouterDefines)),
)
