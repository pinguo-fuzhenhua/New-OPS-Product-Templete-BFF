module github.com/pinguo-icc/Camera360

go 1.17

require (
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20220315030210-f68702a8c2b3
	github.com/go-kratos/kratos/v2 v2.2.0
	github.com/golang/protobuf v1.5.2
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/pinguo-icc/field-definitions v0.0.0-20220630073122-1e7f5984f5db
	github.com/pinguo-icc/go-base/v2 v2.5.13
	github.com/pinguo-icc/go-lib/v2 v2.0.2
	github.com/pinguo-icc/kratos-library/v2 v2.1.8
	github.com/pinguo-icc/operational-basic-svc v0.0.0-20220630071125-878ed2a256c8
	github.com/pinguo-icc/operational-positions-svc v0.0.0-20220517081046-f0bce17829b4
	github.com/pinguo-icc/operations-data-env-svc v0.0.1
	github.com/pinguo-icc/operations-material-svc v0.0.0-20220506062329-901a5474a15f
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/trace v1.4.1
	go.uber.org/zap v1.21.0
	golang.org/x/text v0.3.7
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-redis/redis/v8 v8.11.4 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kisielk/sqlstruct v0.0.0-20210630145711-dae28ed37023 // indirect
	github.com/klauspost/compress v1.14.3 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/olivere/elastic/v7 v7.0.26 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pinguo-icc/bmallsdk v1.0.27 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/qiniu/go-sdk/v7 v7.11.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.8.1 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.10.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.0 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.8.4 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.4.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.4.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.4.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.4.1 // indirect
	go.opentelemetry.io/otel/sdk v1.4.1 // indirect
	go.opentelemetry.io/proto/otlp v0.12.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	google.golang.org/genproto v0.0.0-20220218161850-94dd64e39d7c // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

//replace github.com/pinguo-icc/field-definitions => ../field-definitions-svc
//replace github.com/pinguo-icc/operational-basic-svc => ../operational-basic-svc
