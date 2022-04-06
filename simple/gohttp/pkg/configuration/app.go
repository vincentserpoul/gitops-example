package configuration

type Config struct {
	Application struct {
		Port      int
		PrettyLog bool
		URL       struct {
			Host    string
			Schemes []string
		}
	}
	Observability struct {
		Collector struct {
			Host string
			Port int
		}
	}
}
