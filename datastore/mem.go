package datastore

import (
	"sync"

	_ "../auth"
)

type InMemoryStorage struct {
	lock    sync.RWMutex
	storage map[string]map[string]string
}

// GetFieldsOfRow - returns fields for requested row
func (mem *InMemoryStorage) GetFieldsOfRow(id string, fields []string) (map[string]string, bool) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()
	value, ok := mem.storage[id]
	if ok {
		return value, true
	}
	return nil, false
}

func (mem *InMemoryStorage) SetFieldsOfRow(id string, fields map[string]string) {
	mem.lock.Lock()
	defer mem.lock.Unlock()
	mem.storage[id] = fields
}

func CreateMemoryStorage() *InMemoryStorage {

	mem := &InMemoryStorage{}
	mem.storage = make(map[string]map[string]string)

	return mem
}
