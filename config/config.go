package config

const DEFAULT_CONCURRENCY int = 1
const DEFAULT_NETWORK string = "TCP"

type Configurations struct {
}

type REDISHandlerFactoryConfigs struct {
	//include whatever else is required here for redis
	Network        string
	Addr           string
	DB             int
	Username       string
	Password       string
	MaxRetries     int
	AckDeadline    int
	MaxConcurrency int
}
