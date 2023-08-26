package database

import (
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	_ "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidStatementCaptureDataFromOracleDatabase(t *testing.T) {
	connection := new(connectionMock)
	connection.On("Create").Return(nil)
	connection.On("Destroy").Return(nil)
	err := connection.Create()
	assert.Nil(t, err)
	defer connection.Destroy()
	queries := entities.Query{Mapping: map[int]string{1: TEST_SQL_QUERY}}
	sql := SQL{
		Connection: connection,
		Queries:    queries,
	}
	repository := OracleRepository{sql}

	statement := entities.Statement{
		ID:        1,
		Timestamp: "2022-09-01 11:00:00",
	}

	_, err = repository.Get(statement)
	assert.Nil(t, err)
}
