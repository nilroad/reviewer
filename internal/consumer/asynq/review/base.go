package review

import (
	"club/internal/domain"
	"context"
)

type reviewService interface {
	Review(ctx context.Context, webhook domain.MergeRequestWebhook) error
}
type Consumer struct {
	reviewService reviewService
}

func New(reviewService reviewService) *Consumer {
	return &Consumer{
		reviewService: reviewService,
	}
}
