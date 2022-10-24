// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pinguo-icc/kratos-library/v2/trace"
	"github.com/pinguo-icc/template/internal/application"
	"github.com/pinguo-icc/template/internal/application/v1"
	"github.com/pinguo-icc/template/internal/domain"
	"github.com/pinguo-icc/template/internal/infrastructure/clientset"
	"github.com/pinguo-icc/template/internal/infrastructure/conf"
	"github.com/pinguo-icc/template/internal/infrastructure/server"
)

// Injectors from wire.go:

func initApp(bootstrap *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	http := bootstrap.Http
	recorder := bootstrap.Recorder
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
	html5Config := bootstrap.HTML5
	activitiesParser := domain.NewActivitiesParser(logger, parserFactory, factory, html5Config)
	operationalPos := &v1.OperationalPos{
		ClientSet: clientSet,
		Logger:    logger,
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
	materialPositionsClient := clientSet.MaterialPositionsClient
	materialPositions := &v1.MaterialPositions{
		MP: materialPositionsClient,
	}
	routerDefines := &application.RouterDefines{
		OPos:    operationalPos,
		DataEnv: dataEnv,
		OpBasic: jsonConfig,
		Mpos:    materialPositions,
	}
	httpServer, cleanup2 := server.NewHttpServer(http, recorder, tracerProvider, logger, routerDefines)
	app := newApp(bootstrap, logger, httpServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
