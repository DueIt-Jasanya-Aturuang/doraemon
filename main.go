package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra/config"
	repository2 "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_usecase"
)

func main() {
	config.LogInit()
	config.EnvInit()

	pgConn := config.NewPgConn()
	// pgConn := &sql.DB{}
	redisConn := config.NewRedisConnection()
	// redisConn := &config.RedisImpl{}

	uow := repository2.NewUnitOfWorkRepoSqlImpl(pgConn)
	userRepo := repository2.NewUserRepoSqlImpl(uow)
	accessRepo := repository2.NewAccessRepoSqlImpl(uow)
	appRepo := repository2.NewAppRepoSqlImpl(uow)
	accountRepo := repository2.NewAccountApiRepoImpl(config.AppAccountApi)
	securityRepo := repository2.NewSecuritySqlRepoImpl(uow)
	oauth2Repo := repository2.NewGoogleOauthRepoImpl(config.OauthClientId, config.OauthClientSecret, config.OauthClientRedirectURI)

	userUsecase := _usecase.NewUserUsecaseImpl(userRepo, redisConn)
	authUsecase := _usecase.NewAuthUsecaseImpl(userRepo, accessRepo, accountRepo, securityRepo)
	appUsecase := _usecase.NewAppUsecaseImpl(appRepo)
	oauth2Usecase := _usecase.NewOauth2UsecaseImpl(userRepo, oauth2Repo)
	otpUsecase := _usecase.NewOTPUsecaseImpl(userRepo, redisConn)
	securityUsecase := _usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userHandler := rest.NewUserHandlerImpl(userUsecase, appUsecase, otpUsecase)
	oauth2Handler := rest.NewOauth2HandlerImpl(oauth2Usecase, authUsecase, appUsecase)
	authHandler := rest.NewAuthHandlerImpl(authUsecase, appUsecase, otpUsecase)
	otpHandler := rest.NewOTPHandlerImpl(otpUsecase, appUsecase)
	securityHandler := rest.NewSecurityHandlerImpl(securityUsecase, appUsecase)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)

	r.Put("/auth/reset-password", userHandler.ResetPassword)
	r.Put("/auth/activasi-account", userHandler.ActivasiAccount)
	r.Post("/auth/forgot-password", userHandler.ForgottenPassword)
	r.Put("/auth/forgot-password", userHandler.ResetForgottenPassword)

	r.Post("/auth/login-google", oauth2Handler.LoginWithGoogle)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/register", authHandler.Register)

	r.Post("/auth/otorisasi", securityHandler.ValidateAccess)
	r.Post("/auth/logout", securityHandler.Logout)

	r.Group(func(r chi.Router) {
		middleware.DeletedClient([]string{"activasi-account"})

		r.Post("/auth/otp", otpHandler.OTPGenerate)
	})
	log.Info().Msgf("Server is running on port %s", config.AppPort)
	err := http.ListenAndServe(config.AppPort, r)
	if err != nil {
		log.Err(err).Msg("failed run server")
		os.Exit(1)
	}
}
