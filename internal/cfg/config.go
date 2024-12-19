package cfg

import (
	"errors"
	"fmt"
	"jaeger/internal/logging"
	"jaeger/internal/tracing"
)

type Config struct {
	Log     logging.Config `mapstructure:",squash"`
	Tracing tracing.Config `mapstructure:",squash"`
}

func NewConfig(fn ...func()) (Config, error) {
	c := Config{}

	err := loadEnv(&c, fn...)
	if err != nil {
		return Config{}, fmt.Errorf("config unmarshalling error: %w", err)
	}

	err = validate(c)
	if err != nil {
		return Config{}, fmt.Errorf("config validation error: %w", errors.New(handleValidatorError(c, err)))
	}

	return c, nil
}
