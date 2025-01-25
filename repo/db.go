package repo

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"task-golang/config"
	"task-golang/model"
	"time"
)

var RedisClient *redis.Client
var Db *gorm.DB

func InitPostgresDb() {
	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
	//	config.Props.DbHost, config.Props.DbUser, config.Props.DbPass, config.Props.DbName, config.Props.DbPort)
	//
	//var err error
	//Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.Fatalf("Failed to connect to the database: %v", err)
	//}

	// Optional: Configure connection pool settings
	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatalf("Failed to configure connection pool: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)

	log.Println("Database connection successfully established!")
}

func MigrateDb() error {
	log.Println("MigrateDb.start")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Props.DbHost, config.Props.DbUser, config.Props.DbPass, config.Props.DbName, config.Props.DbPort)

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	errAutoMigrate := Db.AutoMigrate(
		&model.User{},
		&model.Permission{},
		&model.Role{},
		&model.UserRole{},
		&model.Role{},
	)

	if errAutoMigrate != nil {
		return fmt.Errorf("error during migration: %w", errAutoMigrate)
	}

	log.Println("Migrations applied successfully!")
	log.Println("MigrateDb.end")
	return nil
}

func InitRedis() *redis.Client {
	// Create a Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Props.RedisHost + ":" + config.Props.RedisPort, // Redis server address
		Password: "",                                                    // No password set
		DB:       0,                                                     // Use default DB
	})

	// Create a context
	ctx := context.Background()

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis!")

	return RedisClient
}
