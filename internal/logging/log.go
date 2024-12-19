package logging

import (
	"log/slog"
	"os"
)

type Config struct {
	// Log format. Default text
	Format string `mapstructure:"log_format" validate:"oneof=text json" env:"LOG_FORMAT"`
	// Log level. Default info
	Level string `mapstructure:"log_level" validate:"oneof=debug info warn error" env:"LOG_LEVEL"`
}

var logLevelMap = map[string]slog.Level{
	"debug": -4,
	"info":  0,
	"warn":  4,
	"error": 8,
}

func NewLogger(conf Config) *slog.Logger {
	if validateLoggingConf(conf) {
		lg := loggerWithConf(conf)
		lg.Info("logger init success")
		return lg
	} else {
		lg := defaultLogger()
		lg.Info("config validation failed, running logger with default values level=info format=text")
		return lg
	}
}

func validateLoggingConf(conf Config) bool {
	if conf.Level != "debug" &&
		conf.Level != "info" &&
		conf.Level != "warn" &&
		conf.Level != "error" {
		return false
	}
	if conf.Format != "text" &&
		conf.Format != "json" {
		return false
	}
	return true
}

func loggerWithConf(conf Config) *slog.Logger {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(logLevelMap[conf.Level])

	if conf.Format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

func defaultLogger() *slog.Logger {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(logLevelMap["info"])

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}
