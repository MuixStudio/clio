package config

import (
	"flag"

	kratosConfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type UserService struct {
	Endpoint string
}
type Config struct {
	Host        string
	Port        uint16
	Env         string
	Redis       RedisConfig
	UserService UserService
}

func Parse(path string) (Config, error) {
	flag.Parse()
	kc := kratosConfig.New(
		kratosConfig.WithSource(
			// form file
			file.NewSource(path),
			// from env
			env.NewSource("CLIO_"),
		))

	// load config
	if err := kc.Load(); err != nil {
		return Config{}, err
	}
	// parse config
	var cfg Config
	if err := kc.Scan(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
