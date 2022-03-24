package configuration

type Config struct {
	Application struct {
		Port      int
		PrettyLog bool
	}
	Database
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
