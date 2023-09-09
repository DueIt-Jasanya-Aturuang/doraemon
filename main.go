package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/usecase"
)

func main() {
	config.LogInit()
	config.EnvInit()

	pgConn := config.NewPgConn()
	// pgConn := &sql.DB{}
	redisConn := config.NewRedisConnection()
	// redisConn := &config.RedisImpl{}

	uow := repository.NewUnitOfWorkRepoSqlImpl(pgConn)
	userRepo := repository.NewUserRepoSqlImpl(uow)
	accessRepo := repository.NewAccessRepoSqlImpl(uow)
	appRepo := repository.NewAppRepoSqlImpl(uow)
	accountRepo := repository.NewAccountApiRepoImpl(config.AppAccountApi)
	securityRepo := repository.NewSecuritySqlRepoImpl(uow)
	oauth2Repo := repository.NewGoogleOauthRepoImpl(config.OauthClientId, config.OauthClientSecret, config.OauthClientRedirectURI)

	userUsecase := usecase.NewUserUsecaseImpl(userRepo, redisConn)
	authUsecase := usecase.NewAuthUsecaseImpl(userRepo, accessRepo, accountRepo, securityRepo)
	appUsecase := usecase.NewAppUsecaseImpl(appRepo)
	oauth2Usecase := usecase.NewOauth2UsecaseImpl(userRepo, oauth2Repo)
	otpUsecase := usecase.NewOTPUsecaseImpl(userRepo, redisConn)
	securityUsecase := usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userHandler := restapi.NewUserHandlerImpl(userUsecase, appUsecase, otpUsecase)
	oauth2Handler := restapi.NewOauth2HandlerImpl(oauth2Usecase, authUsecase, appUsecase)
	authHandler := restapi.NewAuthHandlerImpl(authUsecase, appUsecase, otpUsecase)
	otpHandler := restapi.NewOTPHandlerImpl(otpUsecase, appUsecase)
	securityHandler := restapi.NewSecurityHandlerImpl(securityUsecase, appUsecase)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)

	r.Post("/auth/reset-password", userHandler.ResetPassword)
	r.Put("/auth/activasi-account", userHandler.ActivasiAccount)
	r.Post("/auth/forgot-password", userHandler.ForgottenPassword)
	r.Put("/auth/forgot-password", userHandler.ResetForgottenPassword)

	r.Post("/auth/login-google", oauth2Handler.LoginWithGoogle)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/register", authHandler.Register)

	r.Post("/auth/authorization", securityHandler.ValidateAccess)
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
