package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	gormlogger "gorm.io/gorm/logger"
)

type Config struct {
	RESTPort                        int    `envconfig:"HTTP_PORT" default:"8080"`
	AllowCORS                       bool   `envconfig:"ALLOW_CORS" default:"true"`
	SystemPassword                  string `envconfig:"SYSTEM_PASSWORD" default:""`
	ParticipantRandomPasswordLength int    `envconfig:"PARTICIPANT_RANDOM_PASSWORD_LENGTH" default:"8"`

	MySQLConfig MySQLConfig
	AuthConfig  AuthConfig
}

type AuthConfig struct {
	LoginTokenExpirationDuration time.Duration `envconfig:"LOGIN_TOKEN_EXPIRATION_DURATION" default:"168h"`
	ApplicationName              string        `envconfig:"APPLICATION_NAME" default:"examitsu"`
	SignatureKey                 []byte        `envconfig:"JWT_SIGNATURE_KEY" default:""`
}

type MySQLConfig struct {
	Username string `envconfig:"MYSQL_USER" default:""`
	Password string `envconfig:"MYSQL_PASSWORD" default:""`
	Hostname string `envconfig:"MYSQL_HOST" default:""`
	DBName   string `envconfig:"MYSQL_DATABASE" default:""`

	// config below can have a default value
	Port         string              `envconfig:"MYSQL_PORT" default:"3306"`
	Charset      string              `envconfig:"MYSQL_CHARSET" default:"utf8mb4"`
	ParseTime    string              `envconfig:"MYSQL_PARSETIME" default:"true"`
	Loc          string              `envconfig:"MYSQL_LOC" default:"Local"`
	GORMLogLevel gormlogger.LogLevel `envconfig:"GORM_LOG_LEVEL" default:"4"`
}

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[config] error loading .env file: %s", err.Error())
	}
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return &cfg
}
