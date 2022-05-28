package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sortednet/statuschecker/internal"
	"github.com/sortednet/statuschecker/internal/statuschecker"
	"github.com/sortednet/statuschecker/internal/store"
	"github.com/sortednet/statuschecker/internal/web"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

var (
	configFile  string = "config/config.yaml"
	initialised bool
)

func main() {
	flag.StringVar(&configFile, "config", "config/config.yaml", "config file")
	flag.Parse()

	logConfig := zap.NewProductionEncoderConfig()
	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := zap.NewProduction()

	appConfig, err := appConfig()
	if err != nil {
		logger.Fatal("Cannot load config", zap.Error(err))
	}

	logger, err = configureLogging(appConfig)
	if err != nil {
		logger.Error("Cannot configure logging", zap.Error(err))
	}
	defer logger.Sync()

	db, err := configureDatabaseConnection(appConfig)
	if err != nil {
		logger.Fatal("Database not available", zap.Error(err))
	}

	checker := configureStatusChecker(appConfig, db)
	webServer, err := configureWebController(appConfig, checker)

	pollContext := context.Background()
	checker.StartPolling(pollContext)

	initialised = true

	webServer.Logger.Fatal(webServer.Start(fmt.Sprintf("0.0.0.0:%s", appConfig.WebPort)))
}

func configureStatusChecker(config internal.Config, db *sql.DB) *statuschecker.StatusChecker {
	queries := store.New(db)
	httpClient := &http.Client{
		Timeout: config.HealthCheckTimeout,
	}
	return statuschecker.NewStatusChecker(queries, config.PollInterval, httpClient)
}

func configureLogging(config internal.Config) (*zap.Logger, error) {

	var cfg zap.Config = zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()

	zap.ReplaceGlobals(logger)
	return logger, err
}

func configureDatabaseConnection(config internal.Config) (*sql.DB, error) {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName)

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, fmt.Errorf("Error opening database %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error connecting to database %w", err)
	}

	return db, nil
}

func configureWebController(config internal.Config, checker *statuschecker.StatusChecker) (*echo.Echo, error) {
	swagger, err := web.GetSwagger()
	if err != nil {
		zap.L().Error("Error loading swagger spec", zap.Error(err))
		return nil, err
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	e := echo.New()
	// Log all requests
	e.Use(echomiddleware.Logger()) // TODO - use the zap logger
	e.Use(echoPrometheus.MetricsMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/ready", ready)
	e.GET("/alive", alive)

	controller := web.NewStatusCheckerController(checker)
	web.RegisterHandlers(e, controller)

	return e, nil
}

func appConfig() (config internal.Config, err error) {
	v := viper.New()
	v.SetEnvPrefix("app")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile(configFile)
	v.AutomaticEnv() // read in environment variables that match

	if err = v.ReadInConfig(); err != nil {
		return
	}
	if err = v.Unmarshal(&config); err != nil {
		return
	}
	zap.L().With(zap.Any("config", config)).Info("app config loaded")
	return config, nil

}

func alive(webCtx echo.Context) error {
	return nil
}

func ready(webCtx echo.Context) error {
	if initialised {
		return nil
	}
	webCtx.Response().Status = http.StatusServiceUnavailable
	return fmt.Errorf("Not ready")
}
