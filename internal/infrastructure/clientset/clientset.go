package clientset

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/conf"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/discovery"
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
	conns, err := newConnection(
		logger,
		traceProvider,
		connInfo{addr: c.FieldDef},
		connInfo{addr: c.OperationalPos, clientOpts: []kgrpc.ClientOption{kgrpc.WithTimeout(5 * time.Second)}},
		connInfo{addr: c.Material},
	)
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

type connInfo struct {
	addr       string
	clientOpts []kgrpc.ClientOption
	dialOpts   []grpc.DialOption
}

func newConnection(logger log.Logger, traceProvider trace.TracerProvider, connData ...connInfo) ([]*grpc.ClientConn, error) {
	r := make([]*grpc.ClientConn, len(connData))
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

	for i := range connData {
		dialOpts := []grpc.DialOption{
			grpc.WithDefaultServiceConfig(retryPolicy),
			grpc.WithBackoffMaxDelay(time.Second),
		}
		dialOpts = append(dialOpts, connData[i].dialOpts...)

		clientOpts := []kgrpc.ClientOption{
			kgrpc.WithDiscovery(discovery.NewDNSDiscovery(log.NewHelper(logger))),
			kgrpc.WithEndpoint(strings.Replace(connData[i].addr, "dns:", "discovery:", 1)),
			kgrpc.WithOptions(dialOpts...),
			kgrpc.WithMiddleware(
				recovery.Recovery(recovery.WithLogger(logger)),
				tracing.Client(tracing.WithTracerProvider(traceProvider)),
				logging.Client(logger),
			),
		}
		clientOpts = append(clientOpts, connData[i].clientOpts...)

		conn, err := kgrpc.DialInsecure(context.TODO(), clientOpts...)
		if err != nil {
			return nil, err
		}
		r[i] = conn
	}

	return r, nil
}
