package datastore

import (
	"github.com/chernser/authtogo/auth"
	"github.com/spf13/viper"
)

// CreateVolatileDataStore - is factory method creating volatile stores from configuration
func CreateVolatileDataStore(conf *viper.Viper) auth.Storage {
	storageType := conf.GetString("storage.volatile.type")
	if storageType == "mem" {
		return &InMemoryStorage{}
	}
	return nil
}

// CreateSecretsDataStore - is factory method creating secrets data store from configuration
func CreateSecretsDataStore(conf *viper.Viper) auth.Storage {
	storageType := conf.GetString("storage.volatile.type")
	if storageType == "mem" {
		return &InMemoryStorage{}
	}
	return nil
}
