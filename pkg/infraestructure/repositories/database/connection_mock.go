package database

import (
	"database/sql"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
)

const TEST_SQL_QUERY = "SELECT sensorId, measurement, timestamp FROM data"

type connectionMock struct {
	mock.Mock
}

func (mock *connectionMock) Configure(entities.Database, entities.Application) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *connectionMock) Create() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *connectionMock) Destroy() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *connectionMock) GetClient() *sql.DB {
	db, mockdb, _ := sqlmock.New()
	mockdb.ExpectQuery(TEST_SQL_QUERY).WillReturnRows()
	return db
}
