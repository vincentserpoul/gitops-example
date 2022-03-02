package configuration

type Config struct {
	Application struct {
		Port int
	}
	Observability struct {
		Collector struct {
			Host string
			Port int
		}
	}
	Stream struct {
		Protocol           string
		Host               string
		Port               int
		InsecureSkipVerify bool
	}
}
