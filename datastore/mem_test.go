package datastore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddReadValues(t *testing.T) {

	store := CreateMemoryStorage()
	storedRow, exists := store.Get("1", []string{"field", "field1"})
	assert.False(t, exists)
	assert.Nil(t, storedRow)

	row := make(map[string]interface{})
	row["field"] = "value_field"
	row["field2"] = "value_field2"
	store.Insert("2", row)

	storedRow, exists = store.Get("2", []string{"field", "field2"})
	assert.True(t, exists)
	assert.NotNil(t, storedRow)
	assert.Equal(t, row["field"], storedRow["field"])

	row["field3"] = "value_field3"
	updated := store.Update("2", row)
	assert.True(t, updated)

	storedRow, exists = store.Get("2", []string{"field", "field2"})
	assert.True(t, exists)
	assert.NotNil(t, storedRow)
	assert.Equal(t, row["field3"], storedRow["field3"])
}

func TestLoadMemoryStoreFromFile(t *testing.T) {

	path := "./valid_memory_file.json"
	storage, err := LoadMemoryStorageFromJsonFile(&path)
	assert.Nil(t, err, "Failed to load json %v", err)

	_, exists := storage.Get("userXXX", []string{"password"})
	assert.False(t, exists)

	userApps, exists := storage.Get("user02", []string{"apps"})
	assert.True(t, exists)
	fmt.Println(userApps)
}
