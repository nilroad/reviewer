package config

import (
	"errors"
	"fmt"

	"git.oceantim.com/backend/packages/golang/database/mysql"
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

	mysqlConfig := mysql.Config{
		Host:         loadString("DATABASE_MYSQL_HOST"),
		Port:         loadInt("DATABASE_MYSQL_PORT"),
		Username:     loadString("DATABASE_MYSQL_USER"),
		Password:     loadString("DATABASE_MYSQL_PASSWORD"),
		DatabaseName: loadString("DATABASE_MYSQL_NAME"),
		Timezone:     loadString("TZ"),
	}

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
			MySQL:      mysqlConfig,
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
	}, nil
}
