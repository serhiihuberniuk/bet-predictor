package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/serhiihuberniuk/bet-predictor/config"
	fetcher "github.com/serhiihuberniuk/bet-predictor/fetcher/es_fetcher"
	"github.com/serhiihuberniuk/bet-predictor/repository"
	"github.com/serhiihuberniuk/bet-predictor/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type CommandContext struct {
	Ctx          context.Context
	Service      *service.Service
	DisconnectFn func()
	CancelFn     func()
}

func commandInit() (CommandContext, error) {
	cfg, err := config.ReadConfig(cfgFileFlag)
	if err != nil {
		return CommandContext{
			DisconnectFn: func() {},
			CancelFn:     func() {},
		}, fmt.Errorf("error while reading config: %w", err)
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoUrl))
	if err != nil {
		return CommandContext{
			DisconnectFn: func() {},
			CancelFn:     func() {},
		}, fmt.Errorf("error while creating mongoDB client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err = mongoClient.Connect(ctx); err != nil {
		return CommandContext{
			Ctx:          ctx,
			CancelFn:     cancel,
			DisconnectFn: func() {},
		}, fmt.Errorf("error while connecting to mongoDB: %w", err)
	}

	disconnectFn := func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println(fmt.Println("error while disconnecting: %w", err))
		}
	}

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return CommandContext{
			Ctx:          ctx,
			DisconnectFn: disconnectFn,
			CancelFn:     cancel,
		}, fmt.Errorf("connection with db is failed: %w", err)
	}

	if apiKeyFlag != "" {
		cfg.ApiKey = apiKeyFlag
	}

	f, err := fetcher.NewESFetcher(ctx, cfg.ApiKey)
	if err != nil {
		return CommandContext{
			Ctx:          ctx,
			DisconnectFn: disconnectFn,
			CancelFn:     cancel,
		}, fmt.Errorf("error while creating fetcher: %w", err)
	}

	s := service.New(repository.New(mongoClient.Database(cfg.MongoDbName)), f)

	return CommandContext{
		Ctx:          ctx,
		Service:      s,
		DisconnectFn: disconnectFn,
		CancelFn:     cancel,
	}, nil
}
