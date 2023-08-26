package utils

import (
	"os"
	"path/filepath"

	sqlEntities "github.com/CESARBR/knot-thing-sql/internal/entities"
	"github.com/janael-pinheiro/knot-cloud-sdk-golang/pkg/entities"
	"gopkg.in/yaml.v2"
)

type config interface {
	sqlEntities.Database | sqlEntities.Application | sqlEntities.Query | map[string]entities.Device | entities.IntegrationKNoTConfig | map[int]string
}

func readTextFile(filepathName string) ([]byte, error) {
	fileContent, err := os.ReadFile(filepath.Clean(filepathName))
	return fileContent, err
}

func ConfigurationParser[T config](filepathName string, configEntity T) (T, error) {
	fileContent, err := readTextFile(filepath.Clean(filepathName))
	if err != nil {
		return configEntity, err
	}

	err = yaml.Unmarshal(fileContent, &configEntity)
	return configEntity, err
}
