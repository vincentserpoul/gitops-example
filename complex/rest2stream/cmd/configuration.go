package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type missingEnvConfigError struct {
	env string
}

func (mece missingEnvConfigError) Error() string {
	return fmt.Sprintf("missing config %s", mece.env)
}

type missingBaseConfigError struct{}

func (mbce missingBaseConfigError) Error() string {
	return "missing base config"
}

func getConfig(currEnv string) (*Config, error) {
	// config
	var cfg Config

	if err := cleanenv.ReadConfig("./config/base.yaml", &cfg); err != nil {
		return nil, missingBaseConfigError{}
	}

	if err := cleanenv.ReadConfig(fmt.Sprintf("./config/%s.yaml", currEnv), &cfg); err != nil {
		return &cfg, missingEnvConfigError{env: currEnv}
	}

	return &cfg, nil
}
