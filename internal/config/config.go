package config

import (
	"time"

	"git.oceantim.com/backend/packages/golang/database/mysql"
	"git.oceantim.com/backend/packages/golang/database/redis"
	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/go-sentry"
)

type AppEnv string

const (
	ProductionEnv AppEnv = "production"
	StageEnv      AppEnv = "stage"
	DevelopEnv    AppEnv = "develop"
	LocalEnv      AppEnv = "local"
)

type Config struct {
	AppEnv         AppEnv
	AppDebug       bool
	LogLevel       log.LogLevelStr
	Locale         string
	HealthToken    string
	Database       Database
	HTTP           HTTP
	Sentry         *sentry.Config
	Tz             string
	TrustedProxies []string
	AsynqMon       AsynqMon
	ReviewWorker   AsynqWorker
	Gitlab         Gitlab
}

type Database struct {
	MySQL      mysql.Config
	Redis      redis.Config
	AsynqRedis *redis.Config
}

type AsynqWorker struct {
	AsynqWorkerCount       int
	AsynqJobProcessTimeout time.Duration
	AsynqMaxRetry          int
	AsynqGroupMaxSize      int
	AsynqGroupGracePeriod  time.Duration
	AsynqGroupMaxDelay     time.Duration
}

type AsynqMon struct {
	APIHost    string
	APIPort    int
	AsynqRedis *redis.Config
}

type HTTP struct {
	APIHost string
	APIPort int
}

type BaseSvcConfig struct {
	BaseURL string
	APIKey  string
}
type CardProxySvcConfig struct {
	BaseURL string
	APIKey  string
}

type ThirdPartyProxySvcConfig struct {
	BaseURL string
	APIKey  string
}

type RuleConfig struct {
	AccountRuleConfig
}

type AccountRuleConfig struct {
	UserNameMaxLength int
	UserNameMinLength int
}

type OKCSSync struct {
	CategorySyncTimes      []string
	BrandSyncTimes         []string
	BranchSyncTimes        []string
	ProductSyncTimes       []string
	BranchProductSyncTimes []string
}

type Gitlab struct {
	BaseURL string
	APIKey  string
}
