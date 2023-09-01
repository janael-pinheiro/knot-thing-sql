package application

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestRemoveSQLPlaceholder(t *testing.T) {
	testTimestamp := "2006-01-02 15:04:05"
	expectedTimestamp := "2006-01-02 15`:`04`:`05"
	actualTimestamp := removeSQLPlaceholder(testTimestamp)
	assert.Equal(t, expectedTimestamp, actualTimestamp)
}

func TestKeepSQLPlaceholder(t *testing.T) {
	testTimestamp := "2006-01-02 15:04:05"
	expectedTimestamp := testTimestamp
	actualTimestamp := keepSQLPlaceholder(testTimestamp)
	assert.Equal(t, expectedTimestamp, actualTimestamp)
}

func TestPutSQLPlaceholder(t *testing.T) {
	testTimestamp := "2006-01-02 15`:`04`:`05"
	expectedTimestamp := "2006-01-02 15:04:05"
	actualTimestamp := putSQLPlaceholder(testTimestamp)
	assert.Equal(t, expectedTimestamp, actualTimestamp)
}

func TestConvertDatetimeToString(t *testing.T) {
	datetimeString := "2022-06-15T08:30:00.0Z"
	convertedDatetime, err := convertStringToDatetime(datetimeString)
	assert.Nil(t, err)
	dateString := convertDatetimeToString(convertedDatetime, entities.CosmosDB)
	expectedconvertedDatetime := "2022-06-15 08`:`30`:`00"
	assert.Equal(t, expectedconvertedDatetime, dateString)
}

func TestGetCurrentDatetime(t *testing.T) {
	brasiliaTime := -3
	expectedCurrentHour := time.Now().UTC().Add(time.Duration(brasiliaTime) * time.Hour).Hour()
	currentDatetime := getCurrentDatetime(brasiliaTime)
	actualCurrentHour := currentDatetime.Hour()
	assert.Equal(t, expectedCurrentHour, actualCurrentHour)
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

func TestLagTimeInSec(t *testing.T) {
	datetimeString := "2022-06-15T08:30:00.0Z"
	convertedDatetime, err := convertStringToDatetime(datetimeString)
	assert.Nil(t, err)
	seconds := 10
	expectedDatetime := convertedDatetime.Add(time.Duration((seconds * -1)) * time.Second)
	actualDatetime := lagTimeInSec(convertedDatetime, seconds)
	assert.Equal(t, expectedDatetime, actualDatetime)
}

func TestFormatTimestampToUTC(t *testing.T) {
	timestamp := "2021-12-22 14:24:00"
	expectedTimestamp := fmt.Sprintf("%s.0-0300", strings.Replace(timestamp, " ", "T", -1))
	actualTimestamp := formatTimestampToUTC(timestamp)
	assert.Equal(t, expectedTimestamp, actualTimestamp)
}
