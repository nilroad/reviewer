package app

import (
	log "git.oceantim.com/backend/packages/golang/go-logger"
	zerologLogger "git.oceantim.com/backend/packages/golang/go-logger/zerologger"
	"os"

	"github.com/rs/zerolog"
)

func InitLogger() {
	z := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.SetDefaultLogger(zerologLogger.NewFromZerolog(z))
}
