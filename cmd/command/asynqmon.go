package command

import (
	"club/internal/config"
	"context"
	"fmt"
	"net/http"
	"time"

	"git.oceantim.com/backend/packages/golang/asynqmon"
	"git.oceantim.com/backend/packages/golang/database/redis"
	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/messaging/v2/asynq"
	"github.com/gin-gonic/gin"
	asynqpkg "github.com/hibiken/asynq"
	"github.com/spf13/cobra"
)

type Asynqmon struct{ logger log.Logger }

func (cmd *Asynqmon) Command(ctx context.Context, cfg *config.Config, logger log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "asynqmon",
		Short: "run asynqmon",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.logger = logger
			cmd.main(ctx, cfg.AsynqMon)
		},
	}
}
func (cmd *Asynqmon) main(ctx context.Context, cfg config.AsynqMon) {
	h := newAsyncMonHTTPHandler(cfg.AsynqRedis)

	defer h.Close()

	r := newGinEngine(h)

	if err := Serve(ctx, fmt.Sprintf("%s:%d", cfg.APIHost, cfg.APIPort), r); err != nil {
		cmd.logger.Fatal(err)
	}
}

func newAsyncMonHTTPHandler(cfg *redis.Config) *asynqmon.HTTPHandler {
	client := asynq.NewClient(cfg)

	return asynqmon.New(asynqmon.Options{
		RootPath: "/",
		RedisConnOpt: asynqpkg.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.Database,
		},
	}, client)
}
func newGinEngine(h *asynqmon.HTTPHandler) *gin.Engine {
	r := gin.Default()

	r.Any(h.RootPath()+"/*any", gin.WrapH(h))

	return r
}

const headerTimeout = 10 * time.Second

func Serve(ctx context.Context, address string, engine http.Handler) error {
	srv := &http.Server{
		Addr:              address,
		Handler:           engine,
		ReadHeaderTimeout: headerTimeout,
	}

	log.Info(fmt.Sprintf("monitoring server starting at: %s", address))
	srvError := make(chan error)
	go func() {
		srvError <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info("monitoring server is shutting down")

		return srv.Shutdown(ctx)
	case err := <-srvError:
		return err
	}
}
