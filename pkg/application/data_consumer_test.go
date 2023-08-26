package application

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func TestConvertStringValueToNumeric(t *testing.T) {
	t.Run("should convert string value to integer", func(t *testing.T) {
		value := "10"
		expected := 10
		result, err := convertStringValueToNumeric(value)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should convert string value to float32", func(t *testing.T) {
		value := "10.53"
		expected := float32(10.53)
		result, err := convertStringValueToNumeric(value)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, result)
	})

	t.Run("should convert string value to float64", func(t *testing.T) {
		value := "10.532313"
		expected := 10.532313
		result, err := convertStringValueToNumeric(value)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, toFixed(result.(float64), 6))
	})

	t.Run("should stay string", func(t *testing.T) {
		value := "this is a string"
		expected := "this is a string"
		result, err := convertStringValueToNumeric(value)
		assert.NoError(t, err)
		assert.Exactly(t, expected, result)
	})

	t.Run("should return 0 if value is empty", func(t *testing.T) {
		value := ""
		expected := 0
		result, err := convertStringValueToNumeric(value)
		assert.Error(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestIsEmptyString(t *testing.T) {
	t.Run("should return true if value is empty", func(t *testing.T) {
		value := ""
		result := isEmptyString(value)
		assert.True(t, result)
	})

	t.Run("should return false if value is not empty", func(t *testing.T) {
		value := "not empty"
		result := isEmptyString(value)
		assert.False(t, result)
	})
}
