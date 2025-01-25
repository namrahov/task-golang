package main

import (
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"task-golang/config"
	"task-golang/handler"
	"task-golang/middleware"
	"task-golang/repo"
)

var opts struct {
	Profile string `short:"p" long:"profile" default:"dev" description:"Application run profile"`
}

func main() {
	// This code snippet parses command-line flags into the 'opts' variable and handles any parsing errors by panicking.
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	// This code initializes the logger for the application.
	initLogger()

	// This code initializes environment variables by calling the initEnvVars function.
	initEnvVars()

	// This line of code loads the configuration settings for the application by calling the LoadConfig function from the config package.
	config.LoadConfig()

	// This code responsible for setting or applying the logging level for the application.
	applyLogLevel()

	log.Info("Application is starting with profile: ", opts.Profile)

	// This code attempts to manipulate db
	err = repo.MigrateDb()
	if err != nil {
		log.Fatal(err)
	}

	redisClient := repo.InitRedis()

	// This code initializes a new HTTP router using the Gorilla Mux package
	router := mux.NewRouter()

	// This line of code applies a middleware function called "Recoverer" from the "mid" package to the router.
	router.Use(mid.Recoverer)

	router.Use(middleware.AuthMiddleware(redisClient))
	// sep application-specific handlers by calling the ApplicationHandler function with the router as an argument.
	handler.UserHandler(router)

	log.Info("Starting server at port: ", config.Props.Port)

	// This code starts an HTTP server that listens on the specified port from the configuration properties and uses the provided router to handle incoming requests.
	log.Fatal(http.ListenAndServe(":"+config.Props.Port, router))
}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	if opts.Profile == "default" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func initEnvVars() {
	if godotenv.Load("profiles/default.env") != nil {
		log.Fatal("Error in loading environment variables from: profiles/default.env")
	} else {
		log.Info("Environment variables loaded from: profiles/default.env")
	}

	if opts.Profile != "default" {
		profileFileName := "profiles/" + opts.Profile + ".env"
		if godotenv.Overload(profileFileName) != nil {
			log.Fatal("Error in loading environment variables from: ", profileFileName)
		} else {
			log.Info("Environment variables overloaded from: ", profileFileName)
		}
	}
}

func applyLogLevel() {
	log.SetLevel(config.Props.LogLevel)
}
