package configuration

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type MissingEnvConfigError struct {
	env string
}

func (mece MissingEnvConfigError) Error() string {
	return fmt.Sprintf("missing config %s", mece.env)
}

type MissingBaseConfigError struct{}

func (mbce MissingBaseConfigError) Error() string {
	return "missing base config"
}

func GetConfig(currEnv string) (*Config, error) {
	// config
	var cfg Config

	if err := cleanenv.ReadConfig("./config/base.yaml", &cfg); err != nil {
		return nil, MissingBaseConfigError{}
	}

	if err := cleanenv.ReadConfig(fmt.Sprintf("./config/%s.yaml", currEnv), &cfg); err != nil {
		return &cfg, MissingEnvConfigError{env: currEnv}
	}

	return &cfg, nil
}
