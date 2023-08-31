package main

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/repository"
)

func main() {
	config.LogInit()
	config.EnvInit()

	pgConn := config.NewPgConn()

	uow := repository.NewUnitOfWorkRepoSqlImpl(pgConn)
	_ = repository.NewUserRepoSqlImpl(uow)
	_ = repository.NewAccessRepoSqlImpl(uow)
	_ = repository.NewAppRepoSqlImpl(uow)
	_ = repository.NewAccountApiRepoImpl(config.AppAccountApi)
	_ = repository.NewGoogleOauthRepoImpl(config.OauthClientId, config.OauthClientSecret, config.OauthClientRedirectURI)
}
