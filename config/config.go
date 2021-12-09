package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var (
	ErrFailedToLoadConfigFile = errors.New("failed to load given config file")
)

type Config struct {
	file          string
	enableEnvVars bool
}

type Option func(*Config)

func New(opts ...Option) (*Config, error) {
	c := &Config{
		enableEnvVars: true,
		file:          "",
	}

	c.add(opts...)

	if c.enableEnvVars {
		viper.AutomaticEnv()
	}

	if c.file == "" {
		return c, nil
	}

	err := readInConfig(c.file)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func readInConfig(file string) error {
	viper.SetConfigFile(file)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrFailedToLoadConfigFile)
	}

	return nil
}

func (c *Config) add(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *Config) File() string {
	return c.file
}
