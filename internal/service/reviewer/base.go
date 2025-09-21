package reviewer

import (
	"club/internal/adapter/githost"
	"text/template"

	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/messaging/v2/message"
)

type Service struct {
	reviewQ          message.Publisher
	gitHostClient    githost.Contract
	reviewResultTmpl *template.Template
	promptTmpl       *template.Template
	logger           log.Logger
}

func New(reviewQ message.Publisher, gitHostClient githost.Contract, logger log.Logger, promptTmpl *template.Template, reviewResultTmpl *template.Template) *Service {
	return &Service{
		reviewQ:          reviewQ,
		gitHostClient:    gitHostClient,
		reviewResultTmpl: reviewResultTmpl,
		promptTmpl:       promptTmpl,
		logger:           logger,
	}
}
