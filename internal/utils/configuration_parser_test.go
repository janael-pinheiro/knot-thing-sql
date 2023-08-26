package utils

import (
	"fmt"
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGivenValidApplicationFilepathReturnConfiguration(t *testing.T) {
	expectedPertinentTags := make(map[int]string)
	expectedPertinentTags[1] = "GR-11-TIT-0410TE001-01"
	expectedConfiguration := entities.Application{
		IntervalBetweenRequestInSeconds: 30,
		PertinentTags:                   expectedPertinentTags,
	}

	var applicationConfig entities.Application

	applicationConfiguration, err := ConfigurationParser("../configuration/application_configuration_test.yaml", applicationConfig)
	fmt.Println(applicationConfiguration)
	assert.Nil(t, err)
	assert.Equal(t, expectedConfiguration.IntervalBetweenRequestInSeconds, applicationConfiguration.IntervalBetweenRequestInSeconds)
	assert.Equal(t, expectedConfiguration.PertinentTags[1], applicationConfiguration.PertinentTags[1])
}

func TestGivenInvalidFilepathReturnError(t *testing.T) {
	var applicationConfig entities.Application

	_, err := ConfigurationParser("invalid_filepath.yaml", applicationConfig)
	assert.NotNil(t, err)
}
