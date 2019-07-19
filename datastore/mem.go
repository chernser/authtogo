package datastore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type InMemoryStorage struct {
	lock    sync.RWMutex
	storage map[string]map[string]interface{}
}

// Get returns fields for requested row
func (mem *InMemoryStorage) Get(id string, fields []string) (map[string]interface{}, bool) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()
	value, ok := mem.storage[id]
	if ok {
		return value, true
	}
	return nil, false
}

// Insert records new record to store
func (mem *InMemoryStorage) Insert(id string, values map[string]interface{}) bool {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	if mem.storage[id] != nil {
		return false
	}

	mem.storage[id] = values
	return true
}

// Update records new values
func (mem *InMemoryStorage) Update(id string, values map[string]interface{}) bool {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	if mem.storage[id] == nil {
		return false
	}
	mem.storage[id] = values
	return true
}

// Delete removes record from datastore
func (mem *InMemoryStorage) Delete(id string) bool {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	if mem.storage[id] == nil {
		return false
	}
	delete(mem.storage, id)
	return true
}

func CreateMemoryStorage() *InMemoryStorage {

	mem := &InMemoryStorage{}
	mem.storage = make(map[string]map[string]interface{})

	return mem
}

func LoadMemoryStorageFromJsonFile(filePath *string) (*InMemoryStorage, error) {

	fileToLoad, err := os.Open(*filePath)
	if err != nil {
		return nil, err
	}

	defer fileToLoad.Close()

	fileContent, _ := ioutil.ReadAll(fileToLoad)
	storage := &InMemoryStorage{}
	json.Unmarshal([]byte(fileContent), &storage.storage)

	return storage, nil
}
