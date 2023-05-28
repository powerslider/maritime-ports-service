package configs

import "log"

// InitializeConfig wires all dependencies for the config module.
func InitializeConfig() *Config {
	conf, err := NewConfig()
	if err != nil {
		log.Fatal("error loading application config")
	}

	return conf
}
