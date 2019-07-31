package datastore

import (
	"github.com/chernser/authtogo/auth"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// CreateVolatileDataStore - is factory method creating volatile stores from configuration
func CreateVolatileDataStore(conf *viper.Viper) auth.Storage {
	storageType := conf.GetString("storage.volatile.type")
	if storageType == "mem" {
		return CreateMemoryStorage()
	}
	return nil
}

// CreateSecretsDataStore - is factory method creating secrets data store from configuration
func CreateSecretsDataStore(conf *viper.Viper) auth.Storage {
	storageType := conf.GetString("storage.secrets.type")
	var storage auth.Storage
	if storageType == "mem" {
		conf.SetDefault("storage.secrets.content_file", "")
		if filePath := conf.GetString("storage.secrets.content_file"); filePath != "" {
			log.Info().Msgf("Loading secrets from file %s ", filePath)
			storage, _ = LoadMemoryStorageFromJsonFile(&filePath)
		} else {
			storage = CreateMemoryStorage()
		}
	}
	return storage
}
