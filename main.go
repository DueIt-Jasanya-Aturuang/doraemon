package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_usecase"
)

func main() {
	infra.LogInit()
	infra.EnvInit()

	pgConn := infra.NewPgConn()
	// pgConn := &sql.DB{}
	redisConn := infra.NewRedisConnection()
	// redisConn := &config.RedisImpl{}

	uow := _repository.NewUnitOfWorkRepositoryImpl(pgConn)
	userRepo := _repository.NewUserRepositoryImpl(uow)
	accessRepo := _repository.NewAccessRepositoryImpl(uow)
	appRepo := _repository.NewAppRepositoryImpl(uow)
	apiService := _repository.NewMicroServiceRepositoryImpl(infra.AppAccountApi)
	securityRepo := _repository.NewSecurityRepositoryImpl(uow)
	oauth2Repo := _repository.NewOauth2RepositoryImpl(infra.OauthClientId, infra.OauthClientSecret, infra.OauthClientRedirectURI)

	userUsecase := _usecase.NewUserUsecaseImpl(userRepo, redisConn)
	authUsecase := _usecase.NewAuthUsecaseImpl(userRepo, accessRepo, apiService, securityRepo)
	appUsecase := _usecase.NewAppUsecaseImpl(appRepo)
	oauth2Usecase := _usecase.NewOauth2UsecaseImpl(userRepo, oauth2Repo)
	otpUsecase := _usecase.NewOTPUsecaseImpl(userRepo, redisConn)
	securityUsecase := _usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userHandler := rest.NewUserHandlerImpl(userUsecase, otpUsecase)
	oauth2Handler := rest.NewOauth2HandlerImpl(oauth2Usecase, authUsecase)
	authHandler := rest.NewAuthHandlerImpl(authUsecase, otpUsecase)
	otpHandler := rest.NewOTPHandlerImpl(otpUsecase)
	securityHandler := rest.NewSecurityHandlerImpl(securityUsecase)
	appMiddleware := middleware.NewAppMiddleware(appUsecase)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	middleware.DeletedClient([]string{"activasi-account"})

	r.Put("/auth/change-password", userHandler.ChangePassword)
	r.Put("/auth/activasi-account", userHandler.ActivasiAccount)
	r.Post("/auth/logout", securityHandler.Logout)

	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.CheckAppID)
		r.Post("/auth/forgot-password", userHandler.ForgottenPassword)
		r.Put("/auth/forgot-password", userHandler.ResetForgottenPassword)
		r.Post("/auth/login-google", oauth2Handler.LoginWithGoogle)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/otorisasi", securityHandler.ValidateAccess)
		r.Post("/auth/otp", otpHandler.GenerateOTP)
	})

	log.Info().Msgf("Server is running on port %s", infra.AppPort)
	err := http.ListenAndServe(infra.AppPort, r)
	if err != nil {
		log.Err(err).Msg("failed run server")
		os.Exit(1)
	}
}
