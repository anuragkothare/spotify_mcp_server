package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Spotify SpotifyConfig `mapstructure:"spotify"`
	Logging LoggingConfig `mapstructure:"logging"`
}

type ServerConfig struct {
	Port         int `mapstructure:"port"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
}

type SpotifyConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
	godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Environment variable bindings
	viper.BindEnv("spotify.client_id", "SPOTIFY_CLIENT_ID")
	viper.BindEnv("spotify.client_secret", "SPOTIFY_CLIENT_SECRET")
	viper.BindEnv("spotify.redirect_uri", "SPOTIFY_REDIRECT_URI")
	viper.BindEnv("server.port", "SERVER_PORT")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Helper function to safely display partial environment variable values
func getPartialEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return "(empty)"
	}
	if len(value) < 8 {
		return value[:len(value)/2] + "..."
	}
	return value[:8] + "..."
}

// Helper function to safely display partial string values
func getPartialString(value string) string {
	if value == "" {
		return "(empty)"
	}
	if len(value) < 8 {
		return value[:len(value)/2] + "..."
	}
	return value[:8] + "..."
}
