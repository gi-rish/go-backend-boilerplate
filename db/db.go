package db

import (
	_db "github.com/NewStreetTechnologies/go-backend-common/db"

	"github.com/NewStreetTechnologies/go-backend-boilerplate/config"
	_logger "github.com/NewStreetTechnologies/go-backend-boilerplate/logger"

	"gorm.io/gorm"
)

var (
	logger = _logger.Logger
)

var DB map[string]*gorm.DB

func initDB() {
	logger.Info("Init database......")
	initDBConnection()
}

func initDBConnection() {
	db := &_db.DB{
		Logger: logger,
	}
	connInfo := db.BuildConnectionInfo(config.GetConfig("database.host"), config.GetConfig("database.port"), config.GetConfig("database.username"), config.GetConfig("database.password"), config.GetConfig("database.name"))
	conn, err := db.ConnectToDB(connInfo, config.GetConfigInNumber("database.logger.level"))
	if err != nil {
		logger.Error("Error when connecting to database")
		panic(err)
	}
	if DB == nil {
		DB = make(map[string]*gorm.DB)
	}
	DB[config.GetConfig("database.name")] = conn
}

func GetDB() *gorm.DB {
	if DB == nil {
		initDB()
	}
	return DB[config.GetConfig("database.name")]
}
