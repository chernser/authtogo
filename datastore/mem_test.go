package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddReadValues(t *testing.T) {

	store := CreateMemoryStorage()
	storedRow, exists := store.GetFieldsOfRow("1", []string{"field", "field1"})
	assert.False(t, exists)
	assert.Nil(t, storedRow)

	row := make(map[string]string)
	row["field"] = "value_field"
	row["field2"] = "value_field2"
	store.SetFieldsOfRow("2", row)

	storedRow, exists = store.GetFieldsOfRow("2", []string{"field", "field2"})
	assert.True(t, exists)
	assert.NotNil(t, storedRow)
}
