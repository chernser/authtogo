package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddReadValues(t *testing.T) {

	store := CreateMemoryStorage()
	storedRow, exists := store.Get("1", []string{"field", "field1"})
	assert.False(t, exists)
	assert.Nil(t, storedRow)

	row := make(map[string]string)
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
