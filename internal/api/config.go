package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config holds the application's environment-based configuration,
// including service name, instance ID, and port.
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark_service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	AppPort     string `envconfig:"APP_PORT" default:"8080"`
	Hostname    string `envconfig:"HOSTNAME"`
}

// NewConfig loads configuration values from environment variables and
// returns a populated Config, generating a random InstanceID if one
// isn't provided.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("api", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg, nil
}
