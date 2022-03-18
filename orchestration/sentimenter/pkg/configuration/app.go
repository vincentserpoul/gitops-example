package configuration

type Config struct {
	Application struct {
		Port         int
		LogPresetDev bool
	}
	Observability struct {
		Collector struct {
			Host string
			Port int
		}
	}
}
