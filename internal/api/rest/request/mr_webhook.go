package request

import "club/internal/domain"

type MRWebhook struct {
	MergeRequestIID   uint64 `json:"merge_request_iid" validate:"required"`
	ProjectID         uint64 `json:"project_id" validate:"required"`
	URL               string `json:"url" validate:"required"`
	CurrentCommitHash string `json:"current_commit_hash" validate:"required"`
	NewCommitHash     string `json:"new_commit_hash" validate:"required"`
}

func (r *MRWebhook) PrepareForValidation() {}

func (r *MRWebhook) ToDomain() *domain.MergeRequestWebhook {
	return &domain.MergeRequestWebhook{
		MergeRequestIID:   r.MergeRequestIID,
		ProjectID:         r.ProjectID,
		URL:               r.URL,
		CurrentCommitHash: r.CurrentCommitHash,
		NewCommitHash:     r.NewCommitHash,
	}
}
