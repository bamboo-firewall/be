package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV" json:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS" json:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT" json:"CONTEXT_TIMEOUT"`
	DBName                 string `mapstructure:"DB_NAME" json:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR" json:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR" json:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET" json:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET" json:"REFRESH_TOKEN_SECRET"`
	MongoDbURI             string `mapstructure:"MONGO_URI" json:"MONGO_URI"`
	MongoDbTimeOut         int    `mapstructure:"MONGO_TIMEOUT" json:"MONGO_TIMEOUT"`
	AdminPassword          string `mapstructure:"ADMIN_PASSWORD" json:"ADMIN_PASSWORD"`
	CORSAllowOrigin        string `mapstructure:"CORS_ALLOW_ORIGIN" json:"CORS_ALLOW_ORIGIN"`
	CORSAllowOMethods      string `mapstructure:"CORS_ALLOW_METHODS" json:"CORS_ALLOW_METHODS"`
	AdminAccount           string `mapstructure:"ADMIN_ACCOUNT" json:"ADMIN_ACCOUNT"`
	EmailDomain            string `mapstructure:"EMAIL_DOMAIN" json:"EMAIL_DOMAIN"`
}

func NewEnv(path string) *Env {
	env := Env{}
	viper.AddConfigPath(path)
	viper.SetConfigType("json")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find config file : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
