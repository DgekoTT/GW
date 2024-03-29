package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type tokenOptions struct {
	JWTRefreshTTL      time.Duration `yaml:"token_refresh_ttl" env-required:"true"`
	JWTAccessTTL       time.Duration `yaml:"token_access_ttl" env-required:"true"`
	JWTVerificationTTL time.Duration `yaml:"token_verification_ttl" env-required:"true"`
}

type Config struct {
	HTTPServer  `yaml:"http_server"`
	Env         string         `yaml:"env" env-default:"local"`
	Options     tokenOptions   `yaml:"token_options"`
	RedisClient *RedisConfig   `yaml:"redis"`
	Mailer      string         `yaml:"mailersend_api_key" env-required:"true"`
	Hcaptcha    string         `yaml:"hcaptcha_secret" env-required:"true"`
	Auth        AuthGRPCConfig `yaml:"auth" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type AuthGRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type RedisConfig struct {
	Host     string `yaml:"redis_host" env-required:"true"`
	Port     int    `yaml:"redis_port" env-required:"true"`
	Password string `yaml:"redis_password" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	// проверяем сушествует ли файл
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	fmt.Println(res, 2)
	return res
}
