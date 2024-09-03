package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
	TokenAuth     *jwtauth.JWTAuth
}

type Config struct {
	Environment     string
	DB_DRIVER       string
	DB_HOST         string
	DB_PORT         string
	DB_USER         string
	DB_PASSWORD     string
	DB_NAME         string
	WEB_SERVER_PORT string
	JWT_SECRET      string
	JWT_EXPIRES_IN  int
	TokenAuth       *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	var cfg conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if _, err := os.Stat(".env"); err == nil {
		// O arquivo .env está presente, então leia-o
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// O arquivo .env não está presente, leia apenas as variáveis de ambiente
		viper.AutomaticEnv()
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg.JWTSecret)
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return &cfg, nil
}

func NewConfig() Config {
	if os.Getenv("ENVIRONMENT") == "" {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading env file")
		}
	}

	jwtExpiresIn, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	return Config{
		Environment:     os.Getenv("ENVIRONMENT"),
		DB_DRIVER:       os.Getenv("DB_DRIVER"),
		DB_HOST:         os.Getenv("DB_HOST"),
		DB_PORT:         os.Getenv("DB_PORT"),
		DB_USER:         os.Getenv("DB_USER"),
		DB_PASSWORD:     os.Getenv("DB_PASSWORD"),
		DB_NAME:         os.Getenv("DB_NAME"),
		WEB_SERVER_PORT: os.Getenv("WEB_SERVER_PORT"),
		JWT_SECRET:      os.Getenv("JWT_SECRET"),
		JWT_EXPIRES_IN:  jwtExpiresIn,
		TokenAuth:       jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil),
	}
}
