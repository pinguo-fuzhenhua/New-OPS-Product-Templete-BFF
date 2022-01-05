package clientset

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/conf"
	fdapi "github.com/pinguo-icc/field-definitions/api"
	oppapi "github.com/pinguo-icc/operational-positions-svc/api"
	opmapi "github.com/pinguo-icc/operations-material-svc/api"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// ClientSet gRPC Client Set
type ClientSet struct {
	fdapi.FieldDefinitionsClient
	oppapi.OperationalPositionsClient
	opmapi.CategoryServiceClient
	opmapi.MaterialServiceClient
}

// NewClientSet new gRPC Client Set
func NewClientSet(c *conf.Clientset, logger log.Logger, traceProvider trace.TracerProvider) (*ClientSet, func(), error) {
	conns, err := newConnection(logger, traceProvider, c.FieldDef, c.OperationalPos, c.Material)
	if err != nil {
		return nil, nil, err
	}

	cancelFn := func() {
		for _, v := range conns {
			v.Close()
		}
	}

	h := &ClientSet{
		FieldDefinitionsClient:     fdapi.NewFieldDefinitionsClient(conns[0]),
		OperationalPositionsClient: oppapi.NewOperationalPositionsClient(conns[1]),
		CategoryServiceClient:      opmapi.NewCategoryServiceClient(conns[2]),
		MaterialServiceClient:      opmapi.NewMaterialServiceClient(conns[2]),
	}

	return h, cancelFn, nil
}

func newConnection(logger log.Logger, traceProvider trace.TracerProvider, addr ...string) ([]*grpc.ClientConn, error) {
	r := make([]*grpc.ClientConn, len(addr))
	retryPolicy := `{
	"LB":"` + wrr.Name + `",
	"MethodConfig": [{
		"Name":[{"Service":""}],
		"RetryPolicy": {
			"MaxAttempts": 3,
			"InitialBackoff": ".01s",
			"MaxBackoff": ".01s", 
			"BackoffMultiplier": 1.0,
			"RetryableStatusCodes": [ "UNAVAILABLE" ]
		}
	}],
	"HealthCheckConfig": {"ServiceName": "grpc.health.v1.Health"}
}`

	for i := range addr {
		conn, err := kgrpc.DialInsecure(
			context.TODO(),
			kgrpc.WithEndpoint(addr[i]),
			kgrpc.WithMiddleware(
				recovery.Recovery(
					recovery.WithLogger(logger),
				),
				tracing.Client(tracing.WithTracerProvider(traceProvider)),
				logging.Client(logger),
			),
			kgrpc.WithTimeout(3*time.Second),
			kgrpc.WithOptions(
				grpc.WithDefaultServiceConfig(retryPolicy),
				grpc.WithBackoffMaxDelay(time.Second),
			),
		)
		if err != nil {
			return nil, err
		}
		r[i] = conn
	}
	return r, nil
}
