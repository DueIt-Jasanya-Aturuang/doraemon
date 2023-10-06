package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	repository2 "github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	usecase2 "github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func main() {
	infra.LogInit()
	infra.EnvInit()

	pgConn := infra.NewPgConn()
	// pgConn := &sql.DB{}
	redisConn := infra.NewRedisConnection()
	// redisConn := &config.RedisImpl{}

	uow := repository2.NewUnitOfWorkRepositoryImpl(pgConn)
	userRepo := repository2.NewUserRepositoryImpl(uow)
	accessRepo := repository2.NewAccessRepositoryImpl(uow)
	appRepo := repository2.NewAppRepositoryImpl(uow)
	apiService := repository2.NewMicroServiceRepositoryImpl(infra.AppAccountApi)
	securityRepo := repository2.NewSecurityRepositoryImpl(uow)
	oauth2Repo := repository2.NewOauth2RepositoryImpl(infra.OauthClientId, infra.OauthClientSecret, infra.OauthClientRedirectURI)

	userUsecase := usecase2.NewUserUsecaseImpl(userRepo, redisConn)
	authUsecase := usecase2.NewAuthUsecaseImpl(userRepo, accessRepo, apiService, securityRepo)
	appUsecase := usecase2.NewAppUsecaseImpl(appRepo)
	oauth2Usecase := usecase2.NewOauth2UsecaseImpl(userRepo, oauth2Repo)
	otpUsecase := usecase2.NewOTPUsecaseImpl(userRepo, redisConn)
	securityUsecase := usecase2.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userHandler := rest.NewUserHandlerImpl(userUsecase, otpUsecase)
	oauth2Handler := rest.NewOauth2HandlerImpl(oauth2Usecase, authUsecase)
	authHandler := rest.NewAuthHandlerImpl(authUsecase, otpUsecase)
	otpHandler := rest.NewOTPHandlerImpl(otpUsecase)
	securityHandler := rest.NewSecurityHandlerImpl(securityUsecase)
	appMiddleware := middleware.NewAppMiddleware(appUsecase)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "App-ID", "User-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.CheckApiKey)
	middleware.DeletedClient([]string{"activasi-account"})

	r.Group(func(r chi.Router) {
		r.Use(middleware.SetAuthorization)
		r.Get("/auth/user", userHandler.GetUserByID)
		r.Put("/auth/change-password", userHandler.ChangePassword)
		r.Put("/auth/change-username", userHandler.ChangeUsername)
		r.Put("/auth/activasi-account", userHandler.ActivasiAccount)
		r.Post("/auth/logout", securityHandler.Logout)
	})

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
