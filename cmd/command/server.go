package command

import (
	"club/internal/api/rest"
	"club/internal/config"
	"club/internal/locale"
	"context"
	"fmt"

	"git.oceantim.com/backend/packages/golang/database/mysql"
	"git.oceantim.com/backend/packages/golang/database/mysql/dbmanager"
	"git.oceantim.com/backend/packages/golang/database/redis"
	"git.oceantim.com/backend/packages/golang/essential/translation"
	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/go-logger/ozonelog"
	"github.com/spf13/cobra"
)

type Server struct {
	logger log.Logger
}

func (cmd *Server) Command(ctx context.Context, logger log.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the HTTP server",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.logger = logger
			cmd.main(ctx, cfg)
		},
	}
}

func (cmd *Server) main(ctx context.Context, cfg *config.Config) {
	db, err := mysql.NewClient(ctx, &cfg.Database.MySQL)
	if err != nil {
		cmd.logger.Fatal("failed to connect to mysql database", log.J{"error": err.Error()})

		return
	}
	defer db.Close()

	gormDB, err := mysql.NewGormWithInstance(db, cfg.AppDebug)
	if err != nil {
		cmd.logger.Fatal(err)
	}

	logger := ozonelog.NewStdOutLogger(cfg.LogLevel, "club:api:db-manager")
	_, _ = dbmanager.New(logger, gormDB, cfg.Database.MySQL)

	_, err = translation.New(cfg.Locale, locale.GetFaTranslations())
	if err != nil {
		cmd.logger.Error("failed to create translator", log.J{"error": err.Error()})

		return
	}

	asynqRedisClient, err := redis.NewClient(ctx, *cfg.Database.AsynqRedis)
	if err != nil {
		cmd.logger.Error("failed to connect to redis", log.J{"error": err.Error()})

		return
	}
	defer asynqRedisClient.Close()

	server := rest.New(cfg.AppEnv, cfg.TrustedProxies, cfg.LogLevel)
	if err := server.Serve(ctx, fmt.Sprintf("%s:%d", cfg.HTTP.APIHost, cfg.HTTP.APIPort)); err != nil {
		cmd.logger.Fatal(err)
	}
}
