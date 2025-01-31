package config

import (
	"github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

const RootPath = "/v1"

type Args struct {
	DbHost                string    `arg:"env:TASK_GOLANG_HOST,required"`
	DbPort                string    `arg:"env:TASK_GOLANG_PORT,required"`
	DbName                string    `arg:"env:TASK_GOLANG_NAME,required"`
	DbUser                string    `arg:"env:TASK_GOLANG_USER,required"`
	DbPass                string    `arg:"env:TASK_GOLANG_PASS,required"`
	RedisHost             string    `arg:"env:TASK_GOLANG_REDIS_HOST,required"`
	RedisPort             string    `arg:"env:TASK_GOLANG_REDIS_PORT,required"`
	LogLevel              log.Level `arg:"env:LOG_LEVEL"`
	Port                  string    `arg:"env:PORT,required"`
	Hostname              string    `arg:"env:HOSTNAME,required"`
	SecretKey             string    `arg:"env:SECRET_KEY,required"`
	UserFrom              string    `arg:"env:USER_FROM,required"`
	UserPassword          string    `arg:"env:USER_PASSWORD,required"`
	UserActivationUrl     string    `arg:"env:USER_ACTIVATION_URL,required"`
	JwtSecret             string    `arg:"env:JWT_SECRET,required"`
	TokenLifetime         string    `arg:"env:TOKEN_LIFETIME,required"`
	TokenExtendedLifetime string    `arg:"env:TOKEN_EXTENDED_LIFETIME,required"`
	SwaggerHost           string    `arg:"env:SWAGGER_HOST,required"`
	SwaggerBasePath       string    `arg:"env:SWAGGER_BASEPATH,required"`
	UrlHeader             string    `arg:"env:URL_HEADER,required"`
	MinioBucket           string    `arg:"env:APP_MINIO_BUCKET,required"`
	AttachmentFileMaxSize string    `arg:"env:ATTACHMENT_FILE_MAX_SIZE,required"`
	TaskVideoMaxSize      string    `arg:"env:TASK_VIDEO_MAX_SIZE,required"`
	MinioUrl              string    `arg:"env:APP_MINIO_URL,required"`
	MinioAccessKey        string    `arg:"env:APP_MINIO_ACCESS_KEY,required"`
	MinioSecretKey        string    `arg:"env:APP_MINIO_SECRET_KEY,required"`
	MinioUseSsl           bool      `arg:"env:APP_MINIO_USE_SSL,required"`
}

var Props Args

func LoadConfig() {
	err := arg.Parse(&Props)
	if err != nil {
		return
	}
}
