package configuration

import "time"

type Config struct {
	Application struct {
		Port int
	}
	CatFact struct {
		URL     string
		Timeout time.Duration
	}
	Sentimenter struct {
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
