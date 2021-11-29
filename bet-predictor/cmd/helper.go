package cmd

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/serhiihuberniuk/bet-predictor/config"
	fetcher "github.com/serhiihuberniuk/bet-predictor/fetcher/es_fetcher"
	"github.com/serhiihuberniuk/bet-predictor/repository"
	"github.com/serhiihuberniuk/bet-predictor/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func commandInit(ctx context.Context, wg *sync.WaitGroup) (*service.Service, error) {
	cfg, err := config.ReadConfig(cfgFileFlag)
	if err != nil {
		return nil, fmt.Errorf("error while reading config: %w", err)
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoUrl))
	if err != nil {
		return nil, fmt.Errorf("error while creating mongoDB client: %w", err)
	}

	if err = mongoClient.Connect(ctx); err != nil {
		return nil, fmt.Errorf("error while connecting to mongoDB: %w", err)
	}

	wg.Add(1)
	go func() {
		<-ctx.Done()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println(fmt.Println("error while disconnecting: %w", err))
		}
		log.Println("database is closed")
		wg.Done()
	}()

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("connection with db is failed: %w", err)
	}

	if apiKeyFlag != "" {
		cfg.ApiKey = apiKeyFlag
	}

	f, err := fetcher.NewESFetcher(ctx, cfg.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("error while creating fetcher: %w", err)
	}

	s := service.New(repository.New(mongoClient.Database(cfg.MongoDbName)), f)

	return s, nil
}
