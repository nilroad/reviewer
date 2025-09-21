package mergereq

import (
	"club/internal/domain"
	"context"

	"git.oceantim.com/backend/packages/golang/essential/validation"
)

type reviewerService interface {
	Prepare(ctx context.Context, webhook domain.MergeRequestWebhook) error
}
type Handler struct {
	validator     *validation.Validation
	reviewService reviewerService
}

func New(validator *validation.Validation, reviewService reviewerService) *Handler {
	return &Handler{
		validator:     validator,
		reviewService: reviewService,
	}
}
