package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoUrl        string `yaml:"mongo_url"`
	MongoDbName     string `yaml:"mongo_db_name"`
	CredentialsFile string `yaml:"credentials_file"`
}

func ReadConfig(configFile string) (Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("error while reading config: %w", err)
	}

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("error while decoding yaml: %w", err)
	}

	return config, nil
}
