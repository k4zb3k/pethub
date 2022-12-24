package config

import (
	"encoding/json"
	"github.com/k4zb3k/pethub/internal/models"
	"github.com/k4zb3k/pethub/pkg/logging"
	"io"
	"os"
)

var logger = logging.GetLogger()

func GetConfig() (*models.Config, error) {
	file, err := os.Open("./config/config.json")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var config models.Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &config, err
}
