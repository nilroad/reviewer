package config

import (
	"errors"
	"fmt"

	"git.oceantim.com/backend/packages/golang/database/redis"
	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/go-sentry"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	appEnv := loadString("APP_ENV")
	debug := loadBool("APP_DEBUG")

	ASYNQRedisConf := &redis.Config{
		Host:     loadString("DATABASE_ASYNQ_REDIS_HOST"),
		Port:     loadString("DATABASE_ASYNQ_REDIS_PORT"),
		Password: loadString("DATABASE_ASYNQ_REDIS_PASSWORD"),
		Database: loadInt("DATABASE_ASYNQ_REDIS_DATABASE"),
	}

	return &Config{
		AppEnv:      AppEnv(appEnv),
		AppDebug:    debug,
		LogLevel:    log.LogLevelStr(loadString("LOG_LEVEL")),
		Locale:      loadString("LOCALE"),
		HealthToken: loadString("HEALTH_TOKEN"),
		Database: Database{
			AsynqRedis: ASYNQRedisConf,
		},
		HTTP: HTTP{
			APIHost: loadString("API_HTTP_HOST"),
			APIPort: loadInt("API_HTTP_PORT"),
		},
		AsynqMon: AsynqMon{
			APIHost:    loadString("ASYNQMON_API_HTTP_HOST"),
			APIPort:    loadInt("ASYNQMON_API_HTTP_PORT"),
			AsynqRedis: ASYNQRedisConf,
		},
		Sentry: &sentry.Config{
			Dsn:                loadString("SENTRY_DSN"),
			EnableTracing:      loadBool("SENTRY_ENABLE_TRACING"),
			TracesSampleRate:   loadFloat64("SENTRY_TRACES_SAMPLE_RATE"),
			Active:             loadBool("SENTRY_ACTIVE"),
			Debug:              debug,
			Environment:        appEnv,
			SampleRate:         loadFloat64("SENTRY_SAMPLE_RATE"),
			ProfilesSampleRate: loadFloat64("SENTRY_PROFILES_SAMPLE_RATE"),
		},
		Tz:             loadString("TZ"),
		TrustedProxies: loadStringSlice("TRUSTED_PROXIES"),
		ReviewWorker: AsynqWorker{
			AsynqWorkerCount:       loadInt("ASYNQ_REVIEW_WORKER_COUNT"),
			AsynqMaxRetry:          loadInt("ASYNQ_REVIEW_WORKER_RETRY"),
			AsynqJobProcessTimeout: loadDuration("ASYNQ_REVIEW_WORKER_TIMEOUT"),
		},
		Gitlab: Gitlab{
			BaseURL: loadString("GITLAB_BASE_URL"),
			APIKey:  loadString("GITLAB_API_KEY"),
		},
	}, nil
}
