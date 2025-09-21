package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/middlewares"
	"git.oceantim.com/backend/packages/golang/go-logger/ozonelog"

	"club/internal/config"

	log "git.oceantim.com/backend/packages/golang/go-logger"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/requestid"
	promkit "github.com/go-kit/kit/metrics/prometheus"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine

	histograms map[string]*promkit.Histogram
}

func New(appEnv config.AppEnv, trustedProxies []string, logLevel log.LogLevelStr) *Server {
	if appEnv == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	logger := ozonelog.NewStdOutLogger(logLevel, "club:api:gin-logger")
	ginLogger := middlewares.NewGinLogger(logger, middlewares.GinLoggerConfig{
		MoreLogStatusStart: http.StatusOK,
	})
	addRequestContext := middlewares.NewAddRequestContext()

	err := r.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatal("trusted proxies has error", log.J{
			"error": err.Error(),
		})
	}

	r.RedirectTrailingSlash = false

	r.Use(
		requestid.New(
			requestid.WithGenerator(func() string {
				return ""
			}),
			requestid.WithCustomHeaderStrKey("x-request-id"),
		),
	)
	r.Use(middlewares.NewTrimmer().Trim(), middlewares.NewConvertToNull().ConvertToNull())

	r.Use(ginLogger.Handle)
	r.Use(gin.RecoveryWithWriter(&middlewares.GinRecoveryWriter{}))
	r.Use(addRequestContext.AddRequestContext())

	r.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	return &Server{
		engine: r,

		histograms: map[string]*promkit.Histogram{},
	}
}

const headerTimeout = 10 * time.Second

func (s *Server) Serve(ctx context.Context, address string) error {
	srv := &http.Server{
		Addr:              address,
		Handler:           s.engine,
		ReadHeaderTimeout: headerTimeout,
	}

	log.Info(fmt.Sprintf("rest server starting at: %s", address))
	srvError := make(chan error)
	go func() {
		srvError <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info("rest server is shutting down")

		return srv.Shutdown(ctx)
	case err := <-srvError:
		return err
	}
}
