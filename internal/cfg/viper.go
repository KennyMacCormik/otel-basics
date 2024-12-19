package cfg

import (
	"github.com/spf13/viper"
)

func setTracingEnv() {
	viper.SetDefault("trace_endpoint", "http://127.0.0.1:4318")
	_ = viper.BindEnv("trace_endpoint")
}

func setLoggingEnv() {
	viper.SetDefault("log_format", "text")
	_ = viper.BindEnv("log_format")

	viper.SetDefault("log_level", "info")
	_ = viper.BindEnv("log_level")
}

func loadEnv(c *Config, regFn ...func()) error {
	setLoggingEnv()
	setTracingEnv()

	for _, fn := range regFn {
		fn()
	}

	viper.AutomaticEnv()
	return viper.Unmarshal(c)
}
