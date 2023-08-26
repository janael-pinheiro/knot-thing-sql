package database

import (
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestListData(t *testing.T) {
	fakeDatabaseMemory := FakeDatabaseMemory{}
	_, err := fakeDatabaseMemory.List()
	assert.Nil(t, err)
}

func TestGetData(t *testing.T) {
	statement := entities.Statement{}
	fakeDatabaseMemory := FakeDatabaseMemory{}
	_, err := fakeDatabaseMemory.Get(statement)
	assert.Nil(t, err)
}
