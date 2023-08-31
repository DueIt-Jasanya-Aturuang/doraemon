package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func EnvInit() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Err(err).Msg("cannot load env file")
		os.Exit(1)
	}

	AppPort = os.Getenv("APPLICATION_PORT")
	AppStatus = os.Getenv("APPLICATION_STATUS")

	PgHost = os.Getenv("DB_POSTGRESQL_HOST")
	PgPort = os.Getenv("DB_POSTGRESQL_PORT")
	PgUser = os.Getenv("DB_POSTGRESQL_USER")
	PgPass = os.Getenv("DB_POSTGRESQL_PASS")
	PgName = os.Getenv("DB_POSTGRESQL_NAME")
	PgSchema = os.Getenv("DB_POSTGRESQL_SCHEMA")
	PgSSL = os.Getenv("DB_POSTGRESQL_SSL")

	dbInt, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisDB = dbInt
	RedisPass = os.Getenv("REDIS_PASS")

	KafkaProtocol = os.Getenv("KAFKA_PROTOCOL")
	KafkaBroker = os.Getenv("KAFKA_BROKER")
	KafkaTopic = os.Getenv("KAFKA_TOPIC")

	DefaultImage = os.Getenv("DEFAULT_DEFAULT_IMAGE")
	AesCFB = os.Getenv("DEFAULT_AES_CFB_KEY")
	AesCBC = os.Getenv("DEFAULT_AES_CBC_KEY")
	AesCBCIV = os.Getenv("DEFAULT_AES_CBC_IV_KEY")

	DefaultKey = os.Getenv("AUTH_DEFAULT_KEY_TOKEN")
	AccessTokenKeyHS = os.Getenv("AUTH_JWT_TOKEN_HS_ACCESS_TOKEN_KEY")
	AccessTokenKeyExpHS = os.Getenv("AUTH_JWT_TOKEN_HS_ACCESS_TOKEN_EXPIRED")
	RefreshTokenKeyHS = os.Getenv("AUTH_JWT_TOKEN_HS_REFRESH_TOKEN_KEY")
	RefreshTokenKeyExpHS = os.Getenv("AUTH_JWT_TOKEN_HS_REFRESH_TOKEN_EXPIRED")

	OauthClientId = os.Getenv("AUTH_OAUTH_GOOGLE_WEB_CLIENT_ID")
	OauthClientSecret = os.Getenv("AUTH_OAUTH_GOOGLE_WEB_CLIENT_SECRET")
	OauthClientRedirectURI = os.Getenv("AUTH_OAUTH_GOOGLE_WEB_REDIRECT_URL")

	log.Info().Msg("config initialization successfully")
}

var (
	AppPort   string
	AppStatus string

	PgHost   string
	PgPort   string
	PgUser   string
	PgPass   string
	PgName   string
	PgSSL    string
	PgSchema string

	RedisHost string
	RedisPort string
	RedisDB   int
	RedisPass string

	KafkaProtocol string
	KafkaBroker   string
	KafkaTopic    string

	DefaultImage string

	AesCFB   string
	AesCBC   string
	AesCBCIV string

	DefaultKey           string
	AccessTokenKeyHS     string
	AccessTokenKeyExpHS  string
	RefreshTokenKeyHS    string
	RefreshTokenKeyExpHS string

	OauthClientId          string
	OauthClientSecret      string
	OauthClientRedirectURI string
)
