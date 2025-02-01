package domain

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	logger "github.com/rismapa/go-banking-lib/config"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`

	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructur:"host"`
		API  string `mapstructure:"apikey"`
	} `mapstructure:"server"`

	DB struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Database string `mapstructure:"name"`
	} `mapstructure:"database"`
}

func GetConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) GetDatabaseConfig() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Database,
	)
}

/*
 * Implemtasi database dengan config dari .env
 */
func (c *Config) GetDatabaseENVConfig() string {
	err := godotenv.Load(".env")
	if err != nil {
		logger.GetLog().Fatal().Err(err).Msg("Error loading .env file")
	}

	dbName := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
}
