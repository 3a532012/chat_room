package conf

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var ENV = os.Getenv("ENV")
var loadConfigOnce sync.Once

type Config struct {
	ENV string

	// main server
	HOST string
	PORT int
	URL  string

	// jwt
	JWT_ISSUER string
	JWT_SECRET string

	// db
	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_DATABASE string
	DB_CONN_STR string

	// micro service
	SINGLE_CHATROOM_SERVICE_HOST string
	SINGLE_CHATROOM_SERVICE_PORT int

	ONLINE_LIST_SERVICE_HOST string
	ONLINE_LIST_SERVICE_PORT int

	USER_EVENT_HOST    string
	USER_EVENT_PORT    int
	USER_EVENT_ADDRESS string

	// micro service event name
	LOGIN_EVENT                 string
	LOGOUT_EVENT                string
	USER_CREATE_EVENT           string
	SENT_MESSAGE_PERSONAL_EVENT string
	READ_MESSAGE_EVENT          string
}

var config Config

func Conf() Config {
	loadConfigOnce.Do(loadConfig)
	return config
}

func dev(config *Config) {
	log.Println("Using dev environment, loading config...")
	config.DB_HOST = "localhost"
	config.DB_PORT = 5432
	config.DB_USER = "dev_user"
	config.DB_PASSWORD = "dev_password"
	config.DB_DATABASE = "dev_database"

	config.JWT_ISSUER = "dev_issuer"
	config.JWT_SECRET = "dev_secret"

	config.SINGLE_CHATROOM_SERVICE_HOST = "localhost"
	config.SINGLE_CHATROOM_SERVICE_PORT = 50051

	config.ONLINE_LIST_SERVICE_HOST = "localhost"
	config.ONLINE_LIST_SERVICE_PORT = 50052

	config.USER_EVENT_HOST = "localhost"
	config.USER_EVENT_PORT = 20031

	config.LOGIN_EVENT = "LOGIN_EVENT"
	config.LOGOUT_EVENT = "LOGOUT_EVENT"
	config.USER_CREATE_EVENT = "USER_CREATE_EVENT"
	config.SENT_MESSAGE_PERSONAL_EVENT = "SENT_MESSAGE_PERSONAL_EVENT"
	config.READ_MESSAGE_EVENT = "READ_MESSAGE_EVENT"
}

func test(config *Config) {
	log.Println("Using test environment, loading config...")
	panic("config: not implemented")
}

func prod(config *Config) {
	log.Println("Using prod environment, loading config...")
	gin.SetMode(gin.ReleaseMode)
	panic("config: not implemented")
}

func loadConfig() {
	config = Config{
		ENV:  ENV,
		HOST: "localhost",
		PORT: 8000,
	}

	switch ENV {
	case "test":
		test(&config)
	case "prod":
		prod(&config)
	default:
		config.ENV = "dev"
		dev(&config)
	}

	if os.Getenv("HOST") != "" {
		config.HOST = os.Getenv("HOST")
	}
	if os.Getenv("PORT") != "" {
		config.PORT, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	config.URL = fmt.Sprintf("%s:%d", config.HOST, config.PORT)
	config.DB_CONN_STR = fmt.Sprintf("mongodb://%s:%d", config.DB_HOST, config.DB_PORT)
	config.USER_EVENT_ADDRESS = fmt.Sprintf("%s:%d", config.USER_EVENT_HOST, config.USER_EVENT_PORT)
}
