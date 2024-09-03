package configs

import (
	"fmt"
	"os"

	"github.com/go-chi/jwtauth"
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
