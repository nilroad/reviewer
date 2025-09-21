package command

import (
	"club/internal/adapter/githost/gitlab"
	"club/internal/config"
	"club/internal/constant"
	"club/internal/consumer/asynq/review"
	"club/internal/service/reviewer"
	"context"
	"fmt"
	"text/template"

	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/go-logger/ozonelog"
	"git.oceantim.com/backend/packages/golang/messaging/v2/asynq"
	"github.com/spf13/cobra"
	gitlabpkg "gitlab.com/gitlab-org/api/client-go"
)

type Consumer struct {
	logger     log.Logger
	reviewFlag bool
}

func (cmd *Consumer) Command(ctx context.Context, logger log.Logger, cfg *config.Config) *cobra.Command {
	consumerCMD := &cobra.Command{
		Use:   "consumer",
		Short: "run club consumer",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.logger = logger
			cmd.main(ctx, cfg)
		},
	}

	consumerCMD.Flags().BoolVar(&cmd.reviewFlag, "review", false, "run review consumer")

	return consumerCMD
}

func (cmd *Consumer) main(ctx context.Context, cfg *config.Config) {
	switch {
	case cmd.reviewFlag:
		cmd.reviewConsumer(ctx, cfg)
	default:
		fmt.Println("No valid command flag was set.")
	}
}

func (cmd *Consumer) reviewConsumer(ctx context.Context, cfg *config.Config) {

	gitlabClient, err := gitlabpkg.NewClient(cfg.Gitlab.APIKey, gitlabpkg.WithBaseURL(cfg.Gitlab.BaseURL))
	if err != nil {
		return
	}

	promptTmpl, err := template.ParseFiles(constant.TemplatePrompt)
	if err != nil {
		cmd.logger.Error(err.Error())
		return
	}

	reviewResultTmpl, err := template.ParseFiles(constant.TemplateReviewResult)
	if err != nil {
		cmd.logger.Error(err.Error())
		return
	}

	gitlabAdapter := gitlab.New(gitlabClient)
	logger := ozonelog.NewStdOutLogger(cfg.LogLevel, "reviewer:consumer:review-service")
	reviewerSvc := reviewer.New(nil, gitlabAdapter, logger, promptTmpl, reviewResultTmpl)

	as := asynq.NewServer(asynq.ServerConfig{
		QueueName:   constant.AsynqReviewQueue,
		WorkerCount: cfg.ReviewWorker.AsynqWorkerCount,
		RedisConfig: *cfg.Database.AsynqRedis,
		Logger:      logger,
	})

	consumer := review.New(reviewerSvc)

	logger = ozonelog.NewStdOutLogger(cfg.LogLevel, "reviewer:consumer:review-consumer")
	cn := asynq.NewConsumer(logger, as, constant.AsynqReviewQueue, consumer.Handle)

	err = cn.Consume(ctx)
	if err != nil {
		cmd.logger.Error(err.Error())
	}
}
