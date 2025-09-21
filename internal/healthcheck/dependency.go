package healthcheck

import (
	"club/internal/config"
	"context"

	"git.oceantim.com/backend/packages/golang/database/mysql"
	"git.oceantim.com/backend/packages/golang/essential/healthcheck"
	log "git.oceantim.com/backend/packages/golang/go-logger"
)

type dependencyName string

func (a dependencyName) String() string {
	return string(a)
}

type Dependency struct {
	mysqlConfig *mysql.Config
	logger      log.Logger
}

func NewDependency(
	cfg *config.Config,
	logger log.Logger,
) *Dependency {
	return &Dependency{
		mysqlConfig: &cfg.Database.MySQL,
		logger:      logger,
	}
}
func (d *Dependency) Get(_ context.Context) map[string]healthcheck.Checker {

	return map[string]healthcheck.Checker{}
}
