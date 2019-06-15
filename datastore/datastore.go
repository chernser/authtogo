package datastore

import (
	"../auth"
	"github.com/spf13/viper"
)

// CreateVolatileDataStore - initializes volatile storage according to configuration.
func CreateVolatileDataStore(conf *viper.Viper) auth.Storage {
	storageType := conf.GetString("volatile_storage.type")
	if storageType == "mem" {
		return &InMemoryStorage{}
	}
	return nil
}
