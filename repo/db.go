package repo

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"task-golang/config"
	"time"
)

var RedisClient *redis.Client
var Db *gorm.DB

func MigrateDb() error {
	log.Println("MigrateDb.start")

	// Create GORM database connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Props.DbUser, config.Props.DbPass, config.Props.DbHost, config.Props.DbPort, config.Props.DbName)

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatalf("Failed to configure connection pool: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	_, errExec := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if errExec != nil {
		return errExec
	}

	log.Println("Migrations applied successfully!")
	log.Println("MigrateDb.end")
	return nil
}

func BeginTransaction() *gorm.DB {
	return Db.Begin()
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
