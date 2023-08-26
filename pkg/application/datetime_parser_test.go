package application

import (
	"testing"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestConvertDatetimeToString(t *testing.T) {
	datetimeString := "2022-06-15T08:30:00.0Z"
	convertedDatetime, err := convertStringToDatetime(datetimeString)
	assert.Nil(t, err)
	dateString := convertDatetimeToString(convertedDatetime, entities.CosmosDB)
	expectedconvertedDatetime := "2022-06-15 08`:`30`:`00"
	assert.Equal(t, expectedconvertedDatetime, dateString)
}

func TestFormatStringDatetimeToUTCWithZone(t *testing.T) {
	datetimeString := "2022-06-15 08`:`30`:`00"
	convertedDatetime := formatTimestampToUTCWithZone(datetimeString)
	expectedconvertedDatetime := "2022-06-15T08:30:00.0Z"
	assert.Equal(t, expectedconvertedDatetime, convertedDatetime)
}
func TestCompareDatetime(t *testing.T) {
	datetimeString := "2022-06-15 08`:`30`:`00"
	datetimeStringWithTimeZone := formatTimestampToUTCWithZone(datetimeString)
	convertedDatetime, err := convertStringToDatetime(datetimeStringWithTimeZone)
	assert.Nil(t, err)
	isTimeLagged := isTimeLagged(convertedDatetime, 5, 30)
	const expectedTimeLagged = true
	assert.Equal(t, expectedTimeLagged, isTimeLagged)
}

func TestUpdateDatetimeWhenLaggedTime(t *testing.T) {
	datetimeString := "2022-06-15 20`:`30`:`00"
	datetimeStringWithTimeZone := formatTimestampToUTCWithZone(datetimeString)
	convertedDatetime, err := convertStringToDatetime(datetimeStringWithTimeZone)
	assert.Nil(t, err)
	isTimeLagged := isTimeLagged(convertedDatetime, 5, 30)

	if isTimeLagged {
		const laggedHours = 5
		const laggedMinSec = 30
		laggedTime := lagTime(convertedDatetime, laggedHours, laggedMinSec)
		convertedDatetimeString := convertDatetimeToString(laggedTime, entities.CosmosDB)
		expectedDatetimeString := "2022-06-15 15`:`30`:`00"
		assert.Equal(t, expectedDatetimeString, convertedDatetimeString)
	}

}

func TestIsLaggedTime(t *testing.T) {
	timestamp := "2022-07-15 19`:`40`:`00"
	laggedHours := 1
	laggedMinSec := 30
	var datetimeStringWithTimeZone string = formatTimestampToUTCWithZone(timestamp)
	convertedDatetime, err := convertStringToDatetime(datetimeStringWithTimeZone)
	if err == nil {
		islagged := isTimeLagged(convertedDatetime, laggedHours, laggedMinSec)
		assert.Equal(t, islagged, true)
	}
}
