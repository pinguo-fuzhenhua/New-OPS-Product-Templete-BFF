package conf

import (
	"github.com/google/wire"
	"github.com/pinguo-icc/kratos-library/v2/trace"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(new(*Bootstrap), "App", "Http", "Trace", "Clientset", "Qiniu", "Params"),
	trace.NewFactory,
)
