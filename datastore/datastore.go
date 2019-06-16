package datastore

import (
	"../auth"
	"github.com/spf13/viper"
)

// CreateVolatileDataStore - is factory method creating volatile stores from configuration
func CreateVolatileDataStore(conf *viper.Viper) auth.Storage {
	storageType := conf.GetString("volatile_storage.type")
	if storageType == "mem" {
		return &InMemoryStorage{}
	}
	return nil
}
