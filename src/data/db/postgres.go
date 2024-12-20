package db

import (
	"fmt"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var logger = logging.NewLogger(config.GetConfig())
var dbClient *gorm.DB

func InitDb(cfg *config.Config) error {
	var err error
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DbName, cfg.Postgres.SslMode, cfg.Postgres.Timezone)

	dbClient, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		return err
	}

	sqlDB, _ := dbClient.DB()

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime * time.Minute)

	logger.Info(logging.Postgres, logging.Startup, "DB connection established", nil)
	//log.Println("DB connection established")
	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	connection, _ := dbClient.DB()
	connection.Close()
}
