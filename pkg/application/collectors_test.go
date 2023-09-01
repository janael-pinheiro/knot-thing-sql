package application

import (
	"database/sql"
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCollector(t *testing.T) {
	fakeCollector := fakeCollector{}
	data := NewDataStrategy(fakeCollector)
	assert.NotNil(t, data)
}

type databaseRepositoryMock struct {
	mock.Mock
}

func (repository *databaseRepositoryMock) List() ([]string, error) {
	return []string{}, nil
}

func (repository *databaseRepositoryMock) Get(statement entities.Statement) (*sql.Rows, error) {
	args := repository.Called(statement)
	return nil, args.Error(0)
}

func (repository *databaseRepositoryMock) ProcessData(rows *sql.Rows) ([]entities.Row, error) {
	args := repository.Called(rows)
	return nil, args.Error(0)
}

func TestGetDatabaseRepository(t *testing.T) {
	repositoryMock := new(databaseRepositoryMock)
	collector := SQLCollector{databaseRepository: repositoryMock}
	databaseRepository := collector.GetDatabaseRepository()
	assert.NotNil(t, databaseRepository)
}

func TestSetCollector(t *testing.T) {
	fakeCollector := fakeCollector{}
	data := NewDataStrategy(fakeCollector)
	err := data.SetCollectorStrategy(fakeCollector)
	assert.Nil(t, err)
}

func TestCollectCollector(t *testing.T) {
	fakeCollector := fakeCollector{}
	data := NewDataStrategy(fakeCollector)
	data.Collect()
	assert.NotNil(t, data)
}
