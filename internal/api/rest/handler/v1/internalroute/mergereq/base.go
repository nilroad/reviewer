package mergereq

import (
	"context"

	"git.oceantim.com/backend/packages/golang/essential/validation"
)

type reviewerService interface {
	Prepare(ctx context.Context, mrIID, projectID, currentCommitHash, newCommitHash string) error
}
type Handler struct {
	validator     validation.Validation
	reviewService reviewerService
}

func New(validator validation.Validation, reviewService reviewerService) *Handler {
	return &Handler{
		validator:     validator,
		reviewService: reviewService,
	}
}
