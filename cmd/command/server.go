package command

import (
	"club/internal/adapter/githost/gitlab"
	"club/internal/api/rest"
	"club/internal/api/rest/handler/v1/internalroute/mergereq"
	"club/internal/config"
	"club/internal/constant"
	health "club/internal/healthcheck"
	"club/internal/locale"
	"club/internal/service/reviewer"
	"context"
	"fmt"

	"git.oceantim.com/backend/packages/golang/essential/exception"
	"git.oceantim.com/backend/packages/golang/essential/healthcheck"
	essmiddleware "git.oceantim.com/backend/packages/golang/essential/middlewares"
	"git.oceantim.com/backend/packages/golang/essential/translation"
	"git.oceantim.com/backend/packages/golang/essential/validation"
	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/go-logger/ozonelog"
	response_formatter "git.oceantim.com/backend/packages/golang/go-response-formatter"
	"git.oceantim.com/backend/packages/golang/messaging/v2/asynq"
	asynq2 "github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	gitlabpkg "gitlab.com/gitlab-org/api/client-go"
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

	translator, err := translation.New(cfg.Locale, locale.GetFaTranslations())
	if err != nil {
		cmd.logger.Error("failed to create translator", log.J{"error": err.Error()})

		return
	}

	exception.SetTranslator(translator)
	validator := validation.New(translator)

	rf := response_formatter.NewResponseFormatter(ozonelog.NewStdOutLogger(
		cfg.LogLevel,
		"club:pkg:response-formatter",
	), cfg.AppDebug)

	asc := asynq.NewClient(cfg.Database.AsynqRedis)

	reviewQ := asynq.NewPublisher(asc, constant.AsynqReviewQueue, []asynq2.Option{
		asynq2.Timeout(cfg.ReviewWorker.AsynqJobProcessTimeout),
		asynq2.MaxRetry(cfg.ReviewWorker.AsynqMaxRetry),
	})

	gitlabClient, err := gitlabpkg.NewClient(cfg.Gitlab.APIKey, gitlabpkg.WithBaseURL(cfg.Gitlab.BaseURL))
	if err != nil {
		return
	}
	gitlabAdapter := gitlab.New(gitlabClient)

	//services
	logger := ozonelog.NewStdOutLogger(cfg.LogLevel, "reviewer:api:reviewer-service")
	reviewService := reviewer.New(reviewQ, gitlabAdapter, logger, nil, nil)

	mrHandler := mergereq.New(validator, reviewService)

	errRespM := essmiddleware.NewHttpErrorResponse(rf)
	sentryCaptureM := essmiddleware.NewCaptureSentry(ozonelog.NewStdOutLogger(
		cfg.LogLevel,
		"club:api:sentry-capture",
	))

	logger = ozonelog.NewStdOutLogger(cfg.LogLevel, "reviewer:health:service")
	dependency := health.NewDependency(cfg, logger)
	healthManager := healthcheck.NewManager(dependency)
	healthcheckH := healthcheck.NewHealthHandler(healthManager, rf, logger)

	server := rest.New(cfg.AppEnv, cfg.TrustedProxies, cfg.LogLevel)

	server.SetupAPIRoutes(errRespM, sentryCaptureM, healthcheckH, mrHandler)

	if err := server.Serve(ctx, fmt.Sprintf("%s:%d", cfg.HTTP.APIHost, cfg.HTTP.APIPort)); err != nil {
		cmd.logger.Fatal(err)
	}
}
