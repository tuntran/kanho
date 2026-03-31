package config

import (
	"log"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	JWT      JWTConfig      `koanf:"jwt"`
	Storage  StorageConfig  `koanf:"storage"`
}

type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type DatabaseConfig struct {
	URL string `koanf:"url"`
}

type JWTConfig struct {
	Secret     string        `koanf:"secret"`
	AccessTTL  time.Duration `koanf:"access_ttl"`
	RefreshTTL time.Duration `koanf:"refresh_ttl"`
}

type StorageConfig struct {
	Endpoint  string `koanf:"endpoint"`
	Bucket    string `koanf:"bucket"`
	AccessKey string `koanf:"access_key"`
	SecretKey string `koanf:"secret_key"`
	UseSSL    bool   `koanf:"use_ssl"`
	PublicURL string `koanf:"public_url"`
}

func Load(path string) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Printf("config file %s not found, using env only: %v", path, err)
	}

	if err := k.Load(env.Provider("KANHO_", ".", func(s string) string {
		return envToKoanfKey(s, "KANHO_")
	}), nil); err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		JWT: JWTConfig{
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 168 * time.Hour,
		},
		Storage: StorageConfig{
			Bucket: "kanho",
			UseSSL: false,
		},
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func envToKoanfKey(s string, prefix string) string {
	key := s[len(prefix):]
	result := make([]byte, 0, len(key))
	for i := 0; i < len(key); i++ {
		if key[i] == '_' {
			result = append(result, '.')
		} else {
			result = append(result, key[i]+'a'-'A')
		}
	}
	return string(result)
}
