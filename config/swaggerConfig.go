package config

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"
	"task-golang/docs"
)

func InitSwagger(router *mux.Router) {
	// Read environment variables for host and base path
	swaggerHost := Props.SwaggerHost
	swaggerBasePath := Props.SwaggerBasePath
	urlHeader := Props.UrlHeader
	// Set dynamic Swagger options
	httpSwagger.URL(urlHeader + swaggerHost + swaggerBasePath + "/swagger/doc.json") // Set the Swagger endpoint dynamically

	// Serve Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Override if you have environment variables like SWAGGER_HOST, SWAGGER_BASEPATH
	docs.SwaggerInfo.Host = getEnv("SWAGGER_HOST", swaggerHost)
	docs.SwaggerInfo.BasePath = getEnv("SWAGGER_BASEPATH", swaggerBasePath)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
