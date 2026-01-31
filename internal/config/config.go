package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env   string `env:"APP_ENV" env-default:"local"`
	HTTP  HTTP
	DB    Postgres
	Valid Validation
}

type Validation struct {
	MyHost          string   `env:"HTTP_HOSTNAME"`
	MaxURLLen       int      `env:"MAX_URL_LEN" env-default:"3"`
	MinAliasLen     int      `env:"MIN_ALIAS_LEN" env-default:"3"`
	MaxAliasLen     int      `env:"MIN_ALIAS_LEN" env-default:"10"`
	DefaultAliasLen int      `env:"GENERATE_ALIAS_LEN" env-default:"6"`
	ForbiddenNames  []string `env:"FORBIDDEN_NAME" env-default:""`
}

type HTTP struct {
	Addr         string        `env:"HTTP_ADDR" env-default:":8000"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	Admin        string        `env:"ADMIN" env-required:"true"`
	Password     string        `env:"ADMIN" env-required:"true"`
}

type Postgres struct {
	URL string `env:"POSTGRES_URL" env-required:"true"`
}

func MustLoad() *Config {
	if os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "local" {
		_ = godotenv.Load()
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("cannot read env: ", err)
	}

	return &cfg
}
