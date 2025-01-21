package config

import (
	"github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

const RootPath = "/v1"

type Args struct {
	DbHost    string    `arg:"env:TASK_GOLANG_HOST,required"`
	DbPort    string    `arg:"env:TASK_GOLANG_PORT,required"`
	DbName    string    `arg:"env:TASK_GOLANG_NAME,required"`
	DbUser    string    `arg:"env:TASK_GOLANG_USER,required"`
	DbPass    string    `arg:"env:TASK_GOLANG_PASS,required"`
	RedisHost string    `arg:"env:TASK_GOLANG_REDIS_HOST,required"`
	RedisPort string    `arg:"env:TASK_GOLANG_REDIS_PORT,required"`
	LogLevel  log.Level `arg:"env:LOG_LEVEL"`
	Port      string    `arg:"env:PORT,required"`
	Hostname  string    `arg:"env:HOSTNAME,required"`
	SecretKey string    `arg:"env:SECRET_KEY,required"`
}

var Props Args

func LoadConfig() {
	err := arg.Parse(&Props)
	if err != nil {
		return
	}
}
