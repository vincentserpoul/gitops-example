package configuration

import "time"

type Config struct {
	Application struct {
		Port         int
		LogPresetDev bool
	}
	CatFact struct {
		URL     string
		Timeout time.Duration
	}
	Sentimenter struct {
		URL     string
		Timeout time.Duration
	}
	Archiver struct {
		URL     string
		Timeout time.Duration
	}
	Observability struct {
		Collector struct {
			Host string
			Port int
		}
	}
}
