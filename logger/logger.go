package logger

import (
	"github.com/NewStreetTechnologies/go-backend-boilerplate/config"

	_logger "github.com/NewStreetTechnologies/go-backend-common/logger"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {

	log := &_logger.Logger{
		LogLevel: config.GetConfigInNumber("logger.level"),
	}

	Logger = log.InitLogger()
}
