package database

import (
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidStatementCaptureDataFromSQLServerDatabase(t *testing.T) {
	connection := new(connectionMock)
	connection.On("Create").Return(nil)
	connection.On("Destroy").Return(nil)
	connection.Create()
	defer connection.Destroy()
	queries := entities.Query{Mapping: map[int]string{1: TEST_SQL_QUERY}}
	sql := SQL{
		Connection: connection,
		Queries:    queries,
	}
	repository := MSSQLServerRepository{sql}

	statement := entities.Statement{
		ID:        1,
		Timestamp: "",
	}

	_, err := repository.Get(statement)
	assert.Nil(t, err)
}
