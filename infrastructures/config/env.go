package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func EnvInit() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Err(err).Msg("cannot load env file")
		os.Exit(1)
	}

	// AppPort = os.Getenv("APPLICATION_PORT")
	// AppStatus = os.Getenv("APPLICATION_STATUS")
	AppAccountApi = os.Getenv("APPLICATION_ACCOUNT_API")

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

	// KafkaProtocol = os.Getenv("KAFKA_PROTOCOL")
	// KafkaBroker = os.Getenv("KAFKA_BROKER")
	// KafkaTopic = os.Getenv("KAFKA_TOPIC")
	//
	// DefaultImage = os.Getenv("DEFAULT_DEFAULT_IMAGE")
	AesCFB = os.Getenv("DEFAULT_AES_CFB_KEY")
	// AesCBC = os.Getenv("DEFAULT_AES_CBC_KEY")
	// AesCBCIV = os.Getenv("DEFAULT_AES_CBC_IV_KEY")
	//
	DefaultKey = os.Getenv("AUTH_DEFAULT_KEY_TOKEN")

	AccessTokenKeyHS = os.Getenv("AUTH_JWT_TOKEN_HS_ACCESS_TOKEN_KEY")
	accessTokenKeyExp, err := time.ParseDuration(os.Getenv("AUTH_JWT_TOKEN_HS_ACCESS_TOKEN_EXPIRED"))
	if err != nil {
		panic(err)
	}
	AccessTokenKeyExpHS = accessTokenKeyExp

	RefreshTokenKeyHS = os.Getenv("AUTH_JWT_TOKEN_HS_REFRESH_TOKEN_KEY")
	refreshTokenKeyExp, err := time.ParseDuration(os.Getenv("AUTH_JWT_TOKEN_HS_REFRESH_TOKEN_EXPIRED"))
	if err != nil {
		panic(err)
	}
	RefreshTokenKeyExpHS = refreshTokenKeyExp

	rememberMeExp, err := time.ParseDuration(os.Getenv("AUTH_JWT_TOKEN_REMEMBER_ME_EXPIRED"))
	if err != nil {
		panic(err)
	}
	RememberMeTokenExp = rememberMeExp

	forgotPasswordTokenExp, err := time.ParseDuration(os.Getenv("AUTH_JWT_TOKEN_FORGOT_PASSWORD_EXPIRED"))
	if err != nil {
		panic(err)
	}
	ForgotPasswordTokenExp = forgotPasswordTokenExp

	OauthClientId = os.Getenv("AUTH_OAUTH_GOOGLE_WEB_CLIENT_ID")
	OauthClientSecret = os.Getenv("AUTH_OAUTH_GOOGLE_WEB_CLIENT_SECRET")
	OauthClientRedirectURI = os.Getenv("AUTH_OAUTH_GOOGLE_WEB_REDIRECT_URL")

	log.Info().Msg("config initialization successfully")
}

var (
	// AppPort       string
	// AppStatus     string

	AppAccountApi string

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

	// KafkaProtocol string
	// KafkaBroker   string
	// KafkaTopic    string

	// DefaultImage string

	AesCFB string
	// AesCBC   string
	// AesCBCIV string

	DefaultKey             string
	AccessTokenKeyHS       string
	AccessTokenKeyExpHS    time.Duration
	RefreshTokenKeyHS      string
	RefreshTokenKeyExpHS   time.Duration
	RememberMeTokenExp     time.Duration
	ForgotPasswordTokenExp time.Duration

	OauthClientId          string
	OauthClientSecret      string
	OauthClientRedirectURI string
)
