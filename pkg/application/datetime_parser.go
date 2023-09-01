package application

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

const datetimeLayout = "2006-01-02 15:04:05"

func removeSQLPlaceholder(timestamp string) string {
	return strings.Replace(timestamp, ":", "`:`", -1)
}

// Used to maintain compatibility with sql placeholder management,
// even if a specific DBMS does not need it, such as SQL Server.
// In this way, the function only returns the received value, without making any changes.
func keepSQLPlaceholder(timestamp string) string {
	return timestamp
}

func putSQLPlaceholder(timestamp string) string {
	return strings.Replace(timestamp, "`:`", ":", -1)
}

func convertStringToDatetime(datetimeString string) (time.Time, error) {
	convertedTime, err := time.Parse(time.RFC3339, datetimeString)
	return convertedTime, err
}

func convertDatetimeToString(datetime time.Time, context string) string {
	datetimeString := datetime.Format(datetimeLayout)
	sqlPlaceholderHandler := getSQLPlaceholderMapping(context)
	return sqlPlaceholderHandler(datetimeString)

}

func getCurrentDatetime(timeOffSet int) time.Time {
	currentDatetime := time.Now()
	nowInUTC := currentDatetime.UTC()
	currentDatetimeInBrasiliaTime := nowInUTC.Add(time.Duration(timeOffSet) * time.Hour)
	return currentDatetimeInBrasiliaTime
}

func isTimeLagged(convertedDatetime time.Time, laggedHours int, minLagSec int) bool {
	const brasiliaTime = -3
	now := getCurrentDatetime(brasiliaTime)
	laggedTime := lagTime(now, laggedHours, minLagSec)
	return convertedDatetime.Before(laggedTime)
}

func formatTimestampToUTCWithZone(timestamp string) string {
	/*
		Expected format for the timestamp: "2021-12-22 14`:`24`:`00".
	*/
	timestamp = putSQLPlaceholder(timestamp)
	formattedTimestamp := strings.Replace(timestamp, " ", "T", -1)
	formattedTimestamp = fmt.Sprintf("%s.0Z", formattedTimestamp)

	return formattedTimestamp
}

func lagTime(datetime time.Time, hoursToLag int, mininumLagSec int) time.Time {
	hoursToLagFloat := float64(hoursToLag)
	minLagHoursFloat := float64(mininumLagSec) / 3600.0
	if hoursToLagFloat < minLagHoursFloat {
		hoursToLagFloat = minLagHoursFloat
	}
	return datetime.Add(time.Duration((hoursToLagFloat*-1.0)*3600.0) * time.Second)
}

func lagTimeInSec(datetime time.Time, seconds int) time.Time {
	return datetime.Add(time.Duration((seconds * -1)) * time.Second)
}

func setTimeStampToCheckNewData(timeStamp string, laggedHours int, intervalBetweenRequestInSeconds int, context string) string {
	retu := ""
	var laggedTime time.Time
	var datetimeStringWithTimeZone string = formatTimestampToUTCWithZone(timeStamp)
	convertedDatetime, err := convertStringToDatetime(datetimeStringWithTimeZone)
	if err == nil {
		laggedTime = lagTimeInSec(convertedDatetime, intervalBetweenRequestInSeconds)
		if isTimeLagged(convertedDatetime, laggedHours, intervalBetweenRequestInSeconds) {
			laggedTime = lagTime(time.Now(), laggedHours, intervalBetweenRequestInSeconds)
		}
	} else {
		laggedTime = lagTimeInSec(time.Now(), intervalBetweenRequestInSeconds)
	}
	retu = convertDatetimeToString(laggedTime, context)
	return retu
}

func updateLaggedLatestTimestamp(intervalBetweenRequestInSeconds int, pertinentTags map[int]string, latestTimestampPerTag map[int]string, laggedHours int, context string) map[int]string {
	/*
		If the most recent record was processed more than X hours ago,
		the timestamp of the last processing is updated to X hours before the current time.
		This procedure ensures that the record was processed no more than X hours ago.
	*/

	lock.Lock()
	defer lock.Unlock()

	if len(latestTimestampPerTag) == 0 || len(pertinentTags) != len(latestTimestampPerTag) {
		currentTimestampPerTag := make(map[int]string)
		for id := range pertinentTags {
			currentTimestampPerTag[id] = setTimeStampToCheckNewData(time.Now().Format("2006-01-02 15:04:05"), laggedHours, intervalBetweenRequestInSeconds, context)
		}
		return currentTimestampPerTag
	}

	for key, timestamp := range latestTimestampPerTag {
		latestTimestampPerTag[key] = setTimeStampToCheckNewData(timestamp, laggedHours, intervalBetweenRequestInSeconds, context)
	}

	return latestTimestampPerTag
}

func formatTimestampToUTC(timestamp string) string {
	/*
		Expected format for the timestamp: "2021-12-22 14:24:00".
	*/
	formattedTimestamp := strings.Replace(timestamp, " ", "T", -1)
	formattedTimestamp = fmt.Sprintf("%s.0-0300", formattedTimestamp)

	return formattedTimestamp
}
