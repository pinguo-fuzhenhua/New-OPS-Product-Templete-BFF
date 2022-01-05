GOBIN?=go

GOMODCACHE:=$(shell $(GOBIN) env GOMODCACHE)
moduleName:=$(shell head -n 1 go.mod | awk '{print $$2}')
kratosVersion:=$(shell fgrep go-kratos/kratos/v2 go.mod | awk '{print $$2}')

# protoc
protocImport:=-I$(GOMODCACHE)/github.com/go-kratos/kratos@$(kratosVersion) -I./ -Iapi
protocOut:=--go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. 
f:=api/api.proto

.PHONY: run proto build wire help

run:
	bin/app -env dev
proto:
	# make protoc f=api/foo.proto
	protoc $(protocImport) $(protocOut) $(f)
build:
	$(GOBIN) build  -o ./bin/ $(moduleName)/cmd/app
wire:
	wire gen $(moduleName)/cmd/app
