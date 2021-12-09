package config

func WithConfigFile(file string) Option {
	return func(config *Config) {
		config.file = file
	}
}

func WithEnvVars(enable bool) Option {
	return func(config *Config) {
		config.enableEnvVars = enable
	}
}
