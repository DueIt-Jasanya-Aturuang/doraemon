package rapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type Presenter struct {
	authUsecase       usecase.AuthUsecase
	oauth2Usecase     usecase.Oauth2Usecase
	otpUsecase        usecase.OTPUsecase
	apiServiceUsecase usecase.ApiServiceUsecase
	securityUsecase   usecase.SecurityUsecase
	userUsecase       usecase.UserUsecase
}

type Dependency struct {
	AuthUsecase       usecase.AuthUsecase
	Oauth2Usecase     usecase.Oauth2Usecase
	OtpUsecase        usecase.OTPUsecase
	ApiServiceUsecase usecase.ApiServiceUsecase
	SecurityUsecase   usecase.SecurityUsecase
	UserUsecase       usecase.UserUsecase
	AppUsecase        usecase.AppUsecase
}

type PresenterConfig struct {
	Dependency *Dependency
}

func NewPresenter(config PresenterConfig) (*http.Server, error) {
	presenter := &Presenter{
		authUsecase:       config.Dependency.AuthUsecase,
		oauth2Usecase:     config.Dependency.Oauth2Usecase,
		otpUsecase:        config.Dependency.OtpUsecase,
		apiServiceUsecase: config.Dependency.ApiServiceUsecase,
		securityUsecase:   config.Dependency.SecurityUsecase,
		userUsecase:       config.Dependency.UserUsecase,
	}

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "App-ID", "User-ID", "Type", "X-Key", "X-Api-Key", "Profile-ID"},
		ExposedHeaders:   []string{"Authorization", "App-ID", "User-ID", "Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.CheckApiKey)
	middleware.DeletedClient([]string{util.ActivasiAccount})
	appMiddleware := middleware.NewAppMiddleware(config.Dependency.AppUsecase)

	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.CheckAppID)
		r.Post("/auth/forgot-password", presenter.ForgottenPassword)
		r.Put("/auth/forgot-password", presenter.ResetForgottenPassword)
		r.Post("/auth/login-google", presenter.LoginWithGoogle)
		r.Post("/auth/login", presenter.Login)
		r.Post("/auth/register", presenter.Register)
		r.Post("/auth/otorisasi", presenter.ValidateAccess)
		r.Post("/auth/otp", presenter.GenerateOTP)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.SetAuthorization)
		r.Get("/auth/user", presenter.GetUserByID)
		r.Put("/auth/change-password", presenter.ChangePassword)
		r.Put("/auth/change-username", presenter.ChangeUsername)
		r.Put("/auth/activasi-account", presenter.ActivasiAccount)
		r.Post("/auth/logout", presenter.Logout)
	})

	server := &http.Server{
		Addr:              infra.AppPort,
		Handler:           r,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
	}

	return server, nil
}
