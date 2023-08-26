package database

import (
	"testing"

	_ "github.com/AntonioJanael/gocosmos"
	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidStatementCaptureDataFromDatabase(t *testing.T) {
	connection := new(connectionMock)
	connection.On("Create").Return(nil)
	connection.On("Destroy").Return(nil)
	err := connection.Create()
	assert.Nil(t, err)
	defer connection.Destroy()
	queries := entities.Query{Mapping: map[int]string{1: TEST_SQL_QUERY}}
	sql := SQL{
		Connection: connection,
		Queries:    queries}

	repository := CosmosDBRepository{sql}

	statement := entities.Statement{
		ID:        1,
		Timestamp: "2022-09-30 11`:`00`:`00",
	}

	_, err = repository.Get(statement)
	assert.Nil(t, err)
}
