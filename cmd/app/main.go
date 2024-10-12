package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/Sadere/song-depository/docs"
	"github.com/Sadere/song-depository/internal/app"
	"github.com/Sadere/song-depository/internal/config"
	"github.com/Sadere/song-depository/internal/database"
	"github.com/Sadere/song-depository/internal/util"
)

//	@title			Songs Depository API v1
//	@version		1.0
//	@description	This API provides functions to store songs along with its info

//	@host		localhost:8080
//	@BasePath	/

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	// Get executable path
	execFile, err := os.Executable()
	if err != nil {
		log.Fatal("failed to get executable path: ", err)
	}
	path := filepath.Dir(execFile)

	// Init config
	cfg, err := config.NewConfig(path)
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	// Init logger
	logger, err := util.NewZapLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("failed to initialize logger: ", err)
	}

	// Init DB
	logger.Debug("PostgreSQL DSN: ", cfg.PostgresDSN)

	db, err := database.NewConnection("pgx", cfg.PostgresDSN)
	if err != nil {
		logger.Fatal("failed to initialize postgresql db: ", err)
	}

	// Create server instance
	app := app.NewServer(cfg, logger, db)

	// Start server
	err = app.Start()
	if err != nil {
		logger.Fatal("failed to run server: ", err)
	}

	// Shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logger.Infoln("graceful server shutdown ...")
}
