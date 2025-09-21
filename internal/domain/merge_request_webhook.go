package domain

type MergeRequestWebhook struct {
	MergeRequestIID   uint64 `json:"merge_request_iid" validate:"required"`
	ProjectID         uint64 `json:"project_id" validate:"required"`
	URL               string `json:"url" validate:"required"`
	CurrentCommitHash string `json:"current_commit_hash" validate:"required"`
	NewCommitHash     string `json:"new_commit_hash" validate:"required"`
}
