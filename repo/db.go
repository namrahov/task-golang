package repo

import (
	"database/sql"
	"fmt"
	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"task-golang/config"
	"time"
)

var Db *pg.DB

func InitDb() {
	Db = pg.Connect(&pg.Options{
		Addr:        config.Props.DbHost + ":" + config.Props.DbPort,
		User:        config.Props.DbUser,
		Password:    config.Props.DbPass,
		Database:    config.Props.DbName,
		PoolSize:    10,
		DialTimeout: 1 * time.Minute,
		MaxRetries:  2,
		MaxConnAge:  15 * time.Minute,
	})
}

func MigrateDb() error {
	log.Info("MigrateDb.start")

	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable", config.Props.DbName,
		config.Props.DbUser, config.Props.DbPass, config.Props.DbHost, config.Props.DbPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Info("Applied ", n, " migrations")
	log.Info("MigrateDb.end")
	return nil
}
