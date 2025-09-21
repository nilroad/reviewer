package review

import (
	"club/internal/domain"
	"club/internal/event"
	"context"
	"encoding/json"

	"git.oceantim.com/backend/packages/golang/messaging/v2/message"
)

func (r *Consumer) Handle(ctx context.Context, msg message.Message) error {
	var evt event.Review
	err := json.Unmarshal(msg.Data, &evt)
	if err != nil {
		return err
	}

	err = r.reviewService.Review(ctx, domain.MergeRequestWebhook{
		MergeRequestIID:   evt.Data.MrIID,
		ProjectID:         evt.Data.ProjectID,
		URL:               evt.Data.URL,
		CurrentCommitHash: evt.Data.CurrentCommitHash,
		NewCommitHash:     evt.Data.NewCommitHash,
	})
	if err != nil {
		return err
	}

	return nil
}
