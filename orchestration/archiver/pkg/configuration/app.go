package configuration

type Config struct {
	Application struct {
		Port         int
		LogPresetDev bool
	}
	Databases struct {
		Read  Database
		Write Database
	}
	Observability struct {
		Collector struct {
			Host string
			Port int
		}
	}
}

type Database struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}
