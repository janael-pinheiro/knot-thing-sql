package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"

	sqlEntities "github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/CESARBR/knot-thing-sql/internal/utils"
	"github.com/CESARBR/knot-thing-sql/pkg/application"
	"github.com/CESARBR/knot-thing-sql/pkg/infraestructure/repositories/database"
	"github.com/CESARBR/knot-thing-sql/pkg/logging"
	"github.com/janael-pinheiro/knot-cloud-sdk-golang/pkg/entities"
	"github.com/janael-pinheiro/knot-cloud-sdk-golang/pkg/gateways/knot"
)

func main() {
	startPprof()
	applicationConfiguration, deviceConfiguration, knotConfiguration, databaseConfiguration := loadConfiguration()

	log := setupLogger(applicationConfiguration.LogFilepath)
	logger := log.Get("Main")

	transmissionChannel := make(chan sqlEntities.CapturedData, len(applicationConfiguration.PertinentTags))
	builderProperties := application.BuilderProperties{
		ApplicationConfiguration: applicationConfiguration,
		Logger:                   log,
		TransmissionChannel:      transmissionChannel,
	}
	connection := database.NewSQLConnection(databaseConfiguration, applicationConfiguration)
	connection.Create()
	defer connection.Destroy()
	buildersMapping := application.NewBuilderMapping()
	builder := application.NewBuilder(buildersMapping, builderProperties)
	builder.SetConnection(connection)
	builder.SetDatabaseConfiguration(databaseConfiguration)
	director := application.NewDirector(builder)
	dataHandler := director.BuildDataHandler()

	go dataHandler.Collect()
	logger.Println("Application started")

	pipeDevices := make(chan map[string]entities.Device)
	knotIntegration, err := knot.NewKNoTIntegration(pipeDevices, knotConfiguration, logger, deviceConfiguration)
	application.VerifyError(err)
	go application.DataConsumer(transmissionChannel, log.Get("Data consumer"), knotIntegration, pipeDevices)
	waitUntilShutdown()
}

func loadConfiguration() (sqlEntities.Application, map[string]entities.Device, entities.IntegrationKNoTConfig, sqlEntities.Database) {
	applicationConfiguration, err := utils.ConfigurationParser("internal/configuration/application_configuration.yaml", sqlEntities.Application{})
	application.VerifyError(err)
	deviceConfiguration, err := utils.ConfigurationParser("internal/configuration/device_config.yaml", make(map[string]entities.Device))
	application.VerifyError(err)
	knotConfiguration, err := utils.ConfigurationParser("internal/configuration/knot_setup.yaml", entities.IntegrationKNoTConfig{})
	application.VerifyError(err)
	databaseConfiguration, err := utils.ConfigurationParser("internal/configuration/database_configuration.yaml", sqlEntities.Database{})
	application.VerifyError(err)
	return applicationConfiguration, deviceConfiguration, knotConfiguration, databaseConfiguration
}

func waitUntilShutdown() {
	quit := make(chan chan struct{})
	<-quit
}

func setupLogger(logFilepath string) *logging.Logrus {
	var log *logging.Logrus
	file, err := os.OpenFile(filepath.Clean(logFilepath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		log = logging.NewLogrus("info", file)
	} else {
		log = logging.NewLogrus("info", os.Stdout)
	}
	return log
}

func startPprof() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
}
