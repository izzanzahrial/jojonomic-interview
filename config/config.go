package config

type Kafka struct {
	Hosts []string
	Topic string
}

type PostgresConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
