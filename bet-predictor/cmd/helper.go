package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/serhiihuberniuk/bet-predictor/config"
	fetcher "github.com/serhiihuberniuk/bet-predictor/fetcher/es-fetcher"
	"github.com/serhiihuberniuk/bet-predictor/repository"
	"github.com/serhiihuberniuk/bet-predictor/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type credentials struct {
	ApiKey string `json:"api_key"`
}

func commandInit() (context.Context, *fetcher.ESFetcher, *service.Service, func(), error) {
	cfg, err := config.ReadConfig(cfgFileFlag)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error while reading config: %w", err)
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoUrl))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error while creating mongoDB client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	if err = mongoClient.Connect(ctx); err != nil {
		return ctx, nil, nil, cancel, fmt.Errorf("error while connecting to mongoDB: %w", err)
	}

	disconnectFunc := func() {
		cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println(fmt.Println("error while disconnecting: %w", err))
		}
	}

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return ctx, nil, nil, disconnectFunc, fmt.Errorf("connection with db is failed: %w", err)
	}

	s := service.New(repository.New(mongoClient.Database(cfg.MongoDbName)))

	var c credentials
	if apiKeyFlag == "" {
		data, err := os.ReadFile(cfg.CredentialsFile)
		if err != nil {
			return ctx, nil, nil, disconnectFunc, fmt.Errorf("error while openinng file with credentials: %w", err)
		}

		if err := json.Unmarshal(data, &c); err != nil {
			return ctx, nil, nil, disconnectFunc, fmt.Errorf("error while unmarshalling credentials file: %w", err)
		}
	} else {
		c.ApiKey = apiKeyFlag
	}

	f, err := fetcher.NewESFetcher(ctx, c.ApiKey)
	if err != nil {
		return ctx, nil, nil, disconnectFunc, fmt.Errorf("error while creating fetcher: %w", err)
	}

	return ctx, f, s, disconnectFunc, nil
}
