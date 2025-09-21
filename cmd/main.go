package main

import (
	"club/cmd/command"
	"club/internal/app"
	"club/internal/config"
	health "club/internal/healthcheck"
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	mysqlmigrate "git.oceantim.com/backend/packages/golang/database/mysql/command"
	"git.oceantim.com/backend/packages/golang/essential/exception"
	"git.oceantim.com/backend/packages/golang/essential/healthcheck"
	healthcheckCommand "git.oceantim.com/backend/packages/golang/essential/healthcheck/command"
	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/go-logger/ozonelog"
	sentryPKG "git.oceantim.com/backend/packages/golang/go-sentry"
	"github.com/spf13/cobra"
)

func main() {
	const description = "club"
	root := &cobra.Command{
		Short: description,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: false,
		},
	}
	app.InitLogger()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sentry := sentryPKG.NewSentry(cfg.Sentry)
	err = sentry.InitSentry()
	if err != nil {
		log.Fatal(err)
	}

	location, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = location

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	exception.SetErrorCodePrefix("REVIEWER")

	var serverCmd command.Server
	var versionCmd command.Version
	var consumerCMD command.Consumer
	var healthcheckCmd healthcheckCommand.SetupHealthCheck
	var asynqmonCMD command.Asynqmon

	loggerServer := ozonelog.NewStdOutLogger(cfg.LogLevel, "club:server")
	loggerhealthcheck := ozonelog.NewStdOutLogger(cfg.LogLevel, "club:command:healthcheck")
	loggerConsumer := ozonelog.NewStdOutLogger(cfg.LogLevel, "club:command:consumer")
	loggerAsynqmon := ozonelog.NewStdOutLogger(cfg.LogLevel, "club:command:asynqmon")

	dependency := health.NewDependency(cfg, loggerhealthcheck)
	manager := healthcheck.NewManager(dependency)

	root.AddCommand(
		versionCmd.Command(),
		serverCmd.Command(ctx, loggerServer, cfg),
		consumerCMD.Command(ctx, loggerConsumer, cfg),
		asynqmonCMD.Command(ctx, cfg, loggerAsynqmon),
		mysqlmigrate.Migrate{}.Command(ctx, &cfg.Database.MySQL),
		healthcheckCmd.Command(ctx, cfg.LogLevel, manager),
	)

	if err := root.Execute(); err != nil {
		log.Panic(fmt.Sprintf("failed to execute root command: \n%v", err))
	}

}
