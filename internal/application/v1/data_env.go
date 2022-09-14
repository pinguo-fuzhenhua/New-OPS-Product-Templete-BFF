package v1

import (
	"context"
	"strings"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pinguo-icc/April/internal/infrastructure/clientset"
	"github.com/pinguo-icc/April/internal/infrastructure/cparam"
	"github.com/pinguo-icc/kratos-library/v2/trace"
	"github.com/pinguo-icc/operations-data-env-svc/api"
)

// 运营位

type DataEnv struct {
	*clientset.ClientSet

	TracerFactory *trace.Factory
}

func (o *DataEnv) ListEnv(ctx khttp.Context) (interface{}, error) {
	ctxTrace, _, span := o.TracerFactory.Debug(context.Context(ctx), "EnvList")
	defer span.End()
	if o.ignoreEnv(ctx) {
		return new(api.ReplyNormalEnvList), nil
	}

	return o.OperationsDataEnvClient.ListNormalEnv(ctxTrace, &empty.Empty{})
}

func (o *DataEnv) ignoreEnv(ctx khttp.Context) bool {
	cp := cparam.FromContext(ctx)
	if strings.ToLower(cp.AppID) == "April" && strings.ToLower(cp.Platform) == "ios" && cp.AppVersion == "9.9.91" {
		return true
	}

	return false

}
