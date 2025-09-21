package command

import (
	"club/internal/config"
	"context"
	"fmt"

	log "git.oceantim.com/backend/packages/golang/go-logger"
	"github.com/spf13/cobra"
)

type Consumer struct {
	logger log.Logger
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

	return consumerCMD
}

func (cmd *Consumer) main(ctx context.Context, cfg *config.Config) {
	switch {
	default:
		fmt.Println("No valid command flag was set.")
	}
}
