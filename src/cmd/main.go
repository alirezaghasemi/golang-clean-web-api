package main

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/cache"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/db"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/db/migrations"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
)

// @SecurityDefinitions.apiKey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)
	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	err = db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	migrations.Up_1()

	api.InitServer(cfg)
}
