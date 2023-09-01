package application

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewTicker(t *testing.T) {
	seconds := 10
	ticker := NewTicker(seconds)
	assert.NotEmpty(t, ticker)
}

func TestHasDataWhenLengthGreaterThanZeroAndErrorNilThenTrue(t *testing.T) {
	var rows []entities.Row
	row := entities.Row{
		Value:     "",
		Timestamp: "",
	}
	var err error
	rows = append(rows, row)
	hasData := hasData(rows, err)
	assert.True(t, hasData)
}

func TestHasDataWhenLengthEqualZeroOrErrorNotNilThenFalse(t *testing.T) {
	var rows []entities.Row
	var err error
	zeroLengthRowsHasData := hasData(rows, err)
	err = errors.New("No data")
	row := entities.Row{
		Value:     "",
		Timestamp: "",
	}
	rows = append(rows, row)
	errorNotNill := hasData(rows, err)
	assert.False(t, zeroLengthRowsHasData)
	assert.False(t, errorNotNill)
}
