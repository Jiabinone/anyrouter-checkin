package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AES      AESConfig
	Admin    AdminConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Path string
}

type JWTConfig struct {
	Secret string
	Expire time.Duration
}

type AESConfig struct {
	Key string
}

type AdminConfig struct {
	Username string
	Password string
}

var C *Config

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	C = &Config{}
	if err := viper.Unmarshal(C); err != nil {
		return err
	}

	C.JWT.Secret = expandEnv(C.JWT.Secret)
	C.AES.Key = expandEnv(C.AES.Key)

	return nil
}

func expandEnv(s string) string {
	if strings.HasPrefix(s, "${") && strings.Contains(s, ":") {
		s = strings.TrimPrefix(s, "${")
		s = strings.TrimSuffix(s, "}")
		parts := strings.SplitN(s, ":", 2)
		if val := os.Getenv(parts[0]); val != "" {
			return val
		}
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return s
}
