package config

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
type Config struct {
	Redis RedisConfig
}
