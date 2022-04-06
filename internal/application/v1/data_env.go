package v1

import (
	"context"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/clientset"
	"github.com/pinguo-icc/kratos-library/v2/trace"
)

// 运营位

type DataEnv struct {
	*clientset.ClientSet

	TracerFactory *trace.Factory
}

func (o *DataEnv) ListEnv(ctx khttp.Context) (interface{}, error) {
	ctxTrace, _, span := o.TracerFactory.Debug(context.Context(ctx), "EnvList")
	defer span.End()

	return o.OperationsDataEnvClient.ListNormalEnv(ctxTrace, &empty.Empty{})
}
