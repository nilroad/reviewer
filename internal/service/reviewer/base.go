package reviewer

import (
	"club/pkg/cursor"

	log "git.oceantim.com/backend/packages/golang/go-logger"
	"git.oceantim.com/backend/packages/golang/messaging/v2/message"
)

type Service struct {
	reviewQ message.Publisher
	cursor  *cursor.Cursor
	logger  log.Logger
}

func New(reviewQ message.Publisher, cursor *cursor.Cursor) *Service {
	return &Service{reviewQ: reviewQ, cursor: cursor}
}
