package main

import (
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"task-golang/config"
	_ "task-golang/docs"
	"task-golang/handler"
	"task-golang/middleware"
	"task-golang/repo"
)

var opts struct {
	Profile string `short:"p" long:"profile" default:"dev" description:"Application run profile"`
}

// @title           Your API Title
// @version         1.0
// @description     This is a sample API server.
// @termsOfService  http://your.terms.of.service.url

// @contact.name   API Support
// @contact.url    http://www.your-support-url.com
// @contact.email  support@your-email.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9093
// @BasePath  /

// @schemes   http https
func main() {
	// Parse flags into the 'opts' variable and handle errors
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	// Initialize logger and environment variables
	initLogger()
	initEnvVars()

	// Load application configuration
	config.LoadConfig()

	// Apply the configured logging level
	applyLogLevel()

	log.Println("Application is starting with profile:", opts.Profile)

	// Perform database migration
	err = repo.MigrateDb()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Redis client
	redisClient := repo.InitRedis()

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Add middlewares
	router.Use(mid.Recoverer)
	router.Use(middleware.AuthMiddleware(redisClient))

	// Set up routes
	setupRoutes(router)

	log.Println("Starting server at port:", config.Props.Port)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":"+config.Props.Port, router))
}

// Function to define and register routes
func setupRoutes(router *mux.Router) {
	// sep application-specific handlers by calling the ApplicationHandler function with the router as an argument.
	handler.UserHandler(router)

	// Swagger handler
	//router.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger"))))
	// Example with Gorilla Mux:
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
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
