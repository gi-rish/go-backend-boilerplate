package config

import (
	"fmt"
	"strconv"

	conf "github.com/NewStreetTechnologies/go-backend-common/config"
)

var (
	_config map[string]string
)

func init() {
	_config = conf.GetConfig("./config.properties", "./decryption_key")
}

func GetConfig(key string) string {
	return _config[key]
}

func GetConfigInNumber(key string) int {
	n, err := strconv.Atoi(_config[key])
	if err != nil {
		panic(fmt.Sprintf("Get config %s error - %d of type %T: %s", key, n, n, err.Error()))
	}
	return n
}

func GetConfigInBool(key string) bool {
	b, err := strconv.ParseBool(_config[key])
	if err != nil {
		panic(fmt.Sprintf("Get config %s error - %t of type %T: %s", key, b, b, err.Error()))
	}
	return b
}

func GetWholeConfig() map[string]string {
	return _config
}
