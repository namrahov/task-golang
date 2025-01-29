package main

import (
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"task-golang/config"
	_ "task-golang/docs"
	"task-golang/handler"
	"task-golang/initializer"
	"task-golang/middleware"
	"task-golang/rabbitmq"
	"task-golang/repo"
)

var opts struct {
	Profile string `short:"p" long:"profile" default:"local" description:"Application run profile"`
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

// @host           {swaggerHost}  // Dynamic host placeholder
// @BasePath       {swaggerBasePath}  // Dynamic base path placeholder
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization

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
	repo.InitRedis()

	//initialize services
	userService := initializer.InitUserService()
	taskService := initializer.InitTaskService()
	boardService := initializer.InitBoardService()
	fileService := initializer.InitFileService()

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Add middlewares
	router.Use(mid.Recoverer)
	router.Use(middleware.AuthMiddleware(userService))

	// sep application-specific handlers by calling the UserHandler function with the router as an argument.
	handler.UserHandler(router, userService)
	handler.BoardHandler(router, boardService)
	handler.TaskHandler(router, taskService)
	handler.FileHandler(router, fileService)

	// Swagger handler
	config.InitSwagger(router)

	// Run RabbitMQ consumer in a goroutine
	go rabbitmq.InitRabbitMq(taskService)

	log.Println("Starting server at port:", config.Props.Port)
	// Start the HTTP server
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
