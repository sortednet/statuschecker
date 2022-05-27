package internal

import "time"

type Config struct {
	LogLevel           string        `json:"logLevel"`
	PollInterval       time.Duration `json:"pollInterval"`
	HealthCheckTimeout time.Duration `json:"healthCheckTimeout"`
	WebPort            string        `json:"port"`

	DbUser     string `json:"dbUser"`
	DbPassword string `json:"dbPassword"`
	DbHost     string `json:"dbHost"`
	DbPort     string `json:"dbPort"`
	DbName     string `json:"dbName"`
}
