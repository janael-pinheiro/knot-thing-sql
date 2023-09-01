package application

import (
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/database"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilderMapping(t *testing.T) {
	builderMapping := NewBuilderMapping()
	assert.Equal(t, new(glassCosmosDBBuilder), builderMapping[entities.CosmosDB])
	assert.NotEqual(t, new(glassCosmosDBBuilder), builderMapping[entities.Oracle])
}

func TestNewBuilder(t *testing.T) {
	builderMapping := NewBuilderMapping()
	applicationConfiguration := new(entities.Application)
	applicationConfiguration.Context = entities.Oracle
	builderProperties := new(BuilderProperties)
	builderProperties.ApplicationConfiguration = *applicationConfiguration
	builder := NewBuilder(builderMapping, *builderProperties)
	oracle := new(restroomOracleBuilder)
	oracle.setProperties(*builderProperties)
	assert.Equal(t, oracle, builder)
}

func TestSetConnection(t *testing.T) {
	builderMapping := NewBuilderMapping()
	applicationConfiguration := new(entities.Application)
	applicationConfiguration.Context = entities.Oracle
	builderProperties := new(BuilderProperties)
	builderProperties.ApplicationConfiguration = *applicationConfiguration
	connection := new(database.SQLConnection)
	builder := NewBuilder(builderMapping, *builderProperties)
	builder.SetConnection(connection)
	oracle := new(restroomOracleBuilder)
	oracle.connection = connection
	oracle.setProperties(*builderProperties)
	assert.Equal(t, oracle, builder)
}
