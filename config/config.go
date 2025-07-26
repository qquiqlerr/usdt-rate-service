package config

import (
	"flag"
	"os"
)

type Config struct {
	GRPCAddress     string
	DatabaseAddress string
	GrinexAddress   string
	LogLevel        string
}

func Load() *Config {
	cfg := &Config{}

	// Set up command line flags
	grpcAddr := flag.String("grpc-addr", "", "GRPC server address")
	dbAddr := flag.String("db-addr", "", "Database address")
	grinexAddr := flag.String("grinex-addr", "", "Grinex address")
	logLevel := flag.String("log-level", "", "Log level")

	flag.Parse()

	// Load configuration values in order of priority: command line flags, environment variables, default values
	cfg.GRPCAddress = getConfigValue(*grpcAddr, "GRPC_ADDRESS", ":50051")
	cfg.DatabaseAddress = getConfigValue(
		*dbAddr,
		"DATABASE_ADDRESS",
		"postgres://usdt_user:usdt_password@0.0.0.0:5432/usdt_rates?sslmode=disable",
	)
	cfg.GrinexAddress = getConfigValue(*grinexAddr, "GRINEX_ADDRESS", "https://grinex.io")
	cfg.LogLevel = getConfigValue(*logLevel, "LOG_LEVEL", "info")

	return cfg
}

// getConfigValue retrieves the configuration value based on the following priority:
// 1. Command line flag value
// 2. Environment variable value
// 3. Default value.
func getConfigValue(flagValue, envKey, defaultValue string) string {
	// Priority 1: command line flag value
	if flagValue != "" {
		return flagValue
	}

	// Priority 2: environment variable value
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}

	// Priority 3: default value
	return defaultValue
}
