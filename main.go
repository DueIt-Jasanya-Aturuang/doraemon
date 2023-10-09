package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository/access_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository/apiService_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository/app_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository/oauth2_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository/security_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository/uow_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository/user_repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/apiService_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/app_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/auth_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/oauth2_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/otp_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/security_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/user_usecase"
)

func main() {
	infra.LogInit()
	infra.EnvInit()

	pgConn := infra.NewPgConn()
	defer func() {
		if err := pgConn.Close(); err != nil {
			log.Warn().Msgf("failed close connection postgres | err %v", err)
		}
	}()

	redisConn := infra.NewRedisConnection()
	defer func() {
		if err := redisConn.Client.Close(); err != nil {
			log.Warn().Msgf("failed close redis client | err %v", err)
		}
	}()

	depen := dependency(pgConn, redisConn)

	httpServer, err := rapi.NewPresenter(rapi.PresenterConfig{
		Dependency: depen,
	})
	if err != nil {
		log.Fatal().Msgf("creating new presenter: %s", err.Error())
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)
	go func() {
		<-exitSignal
		log.Info().Msg("Interrupt signal received, exiting...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()

		log.Info().Msg("shutting down HTTP server")
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Err(err).Msg("shutting down HTTP server")
		}

	}()

	log.Info().Msgf("Server is running on port %s", infra.AppPort)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msgf("serving HTTP server: %s", err.Error())
	}
}

func dependency(db *sql.DB, cache *infra.RedisImpl) *rapi.Dependency {
	uow := uow_repository.NewUnitOfWorkRepositoryImpl(db)
	userRepo := user_repository.NewUserRepositoryImpl(uow)
	accessRepo := access_repository.NewAccessRepositoryImpl(uow)
	appRepo := app_repository.NewAppRepositoryImpl(uow)
	apiServiceRepo := apiService_repository.NewApiServiceRepositoryImpl()
	securityRepo := security_repository.NewSecurityRepositoryImpl(uow)
	oauth2Repo := oauth2_repository.NewOauth2RepositoryImpl()

	apiServiceUsecase := apiService_usecase.NewApiServiceUsecaseImpl(apiServiceRepo)
	securityUsecase := security_usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)
	authUsecase := auth_usecase.NewAuthUsecaseImpl(userRepo, accessRepo, apiServiceUsecase, securityUsecase)
	oauth2Usecase := oauth2_usecase.NewOauth2UsecaseImpl(userRepo, oauth2Repo)
	otpUsecase := otp_usecase.NewOTPUsecaseImpl(userRepo, cache)
	userUsecase := user_usecase.NewUserUsecaseImpl(userRepo, cache)
	appUsecase := app_usecase.NewAppUsecaseImpl(appRepo)

	return &rapi.Dependency{
		AuthUsecase:       authUsecase,
		Oauth2Usecase:     oauth2Usecase,
		OtpUsecase:        otpUsecase,
		ApiServiceUsecase: apiServiceUsecase,
		SecurityUsecase:   securityUsecase,
		UserUsecase:       userUsecase,
		AppUsecase:        appUsecase,
	}
}
