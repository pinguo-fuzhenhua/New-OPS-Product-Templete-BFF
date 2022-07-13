// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pinguo-icc/Camera360/internal/application"
	"github.com/pinguo-icc/Camera360/internal/application/v1"
	"github.com/pinguo-icc/Camera360/internal/domain"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/clientset"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/conf"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/server"
	"github.com/pinguo-icc/kratos-library/v2/trace"
)

// Injectors from wire.go:

func initApp(bootstrap *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	http := bootstrap.Http
	confLog := bootstrap.Log
	config := bootstrap.Trace
	tracerProvider := trace.NewTracerProvider(config)
	confClientset := bootstrap.Clientset
	clientSet, cleanup, err := clientset.NewClientSet(confClientset, logger, tracerProvider)
	if err != nil {
		return nil, nil, err
	}
	fieldDefinitionsClient := clientSet.FieldDefinitionsClient
	parserFactory := domain.NewParserFactory(fieldDefinitionsClient)
	factory := trace.NewFactory(config)
	activitiesParser := domain.NewActivitiesParser(logger, parserFactory, factory)
	operationalPos := &v1.OperationalPos{
		ClientSet: clientSet,
		Parser:    activitiesParser,
	}
	dataEnv := &v1.DataEnv{
		ClientSet:     clientSet,
		TracerFactory: factory,
	}
	operationalBasicClient := clientSet.OperationalBasicClient
	jsonConfig := &v1.JsonConfig{
		Opbasic: operationalBasicClient,
	}
	routerDefines := &application.RouterDefines{
		OPos:    operationalPos,
		DataEnv: dataEnv,
		OpBasic: jsonConfig,
	}
	httpServer, cleanup2 := server.NewHttpServer(http, confLog, tracerProvider, logger, routerDefines)
	app := newApp(bootstrap, logger, httpServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
