package app

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node/cmd/meta-node-dns/internal/api/handler"
	"github.com/meta-node-blockchain/meta-node/cmd/meta-node-dns/internal/api/route"
	"github.com/meta-node-blockchain/meta-node/cmd/meta-node-dns/internal/config"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
)

type App struct {
	config *config.AppConfig

	dnsHandler *handler.DnsHandler
	dnsStorage storage.Storage

	api *gin.Engine
}

func NewApp(
	configFilePath string,
	logLevel int,
) (*App, error) {
	var loggerConfig = &logger.LoggerConfig{
		Flag:    logLevel,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)

	app := &App{}
	// load config
	var err error
	app.config, err = config.LoadConfig(configFilePath)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		return nil, err
	}

	// Initialize the database
	app.dnsStorage, err = storage.LoadDb(
		app.config.DnsDBPath,
		"",
	)
	if err != nil {
		log.Fatal("Failed to open dns db")
		panic("Failed to open dns db")
	}

	app.dnsHandler = handler.NewDnsHandler(
		app.dnsStorage,
		app.config,
	)

	// Initialize the Gin router
	r := gin.Default()
	// Initialize cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	// Setup the user API routes
	apiRouter := r.Group("/api/")
	route.SetupBlockRoutes(apiRouter, app.dnsHandler) // Initialize user routes
	//
	app.api = r

	return app, nil
}

func (app *App) Run() {
	// Run the server
	err := app.api.Run(app.config.APIAddress)
	if err != nil {
		log.Fatal("Failed to start the server")
		panic("Failed to start the server")
	}
}

func (app *App) Stop() error {
	app.dnsStorage.Close()
	logger.Warn("App Stopped")
	return nil
}
