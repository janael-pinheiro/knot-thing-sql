package application

import (
	"fmt"
	"strconv"
	"sync"

	sqlEntities "github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/janael-pinheiro/knot-cloud-sdk-golang/pkg/entities"
	"github.com/janael-pinheiro/knot-cloud-sdk-golang/pkg/gateways/knot"
	"github.com/sirupsen/logrus"
)

var consumerMutex *sync.Mutex = knot.GetMutex()

func DataConsumer(transmissionChannel chan sqlEntities.CapturedData, logger *logrus.Entry, knotIntegration *knot.Integration, pipeDevices chan map[string]entities.Device) {
	/*
		Receives the data collected from the database.
	*/
	devices := <-pipeDevices
	device := knotIntegration.GetDevice(devices)
	device = knotIntegration.Register(device)
	for capturedData := range transmissionChannel {
		for _, row := range capturedData.Rows {
			var sensors []entities.Data
			convertedValue, err := convertStringValueToNumeric(row.Value)
			if err == nil {
				sensor := entities.Data{SensorID: capturedData.ID, Value: convertedValue, TimeStamp: formatTimestampToUTC(row.Timestamp)}
				sensors = append(sensors, sensor)
			}
			knotIntegration.SentDataToKNoT(sensors, device)
			//Reset the sensors array to avoid data duplication.
			device.Data = nil
			sensors = nil
		}
	}
}

func convertStringValueToNumeric(value string) (interface{}, error) {
	if integerValue, err := strconv.Atoi(value); err == nil {
		return integerValue, err
	}
	if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
		return floatValue, err
	}
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue, err
	}
	if isEmptyString(value) {
		return 0, fmt.Errorf("type conversion error")
	}
	return value, nil
}

func isEmptyString(value string) bool {
	return value == ""
}
