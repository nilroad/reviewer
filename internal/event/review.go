package event

type Review struct {
	Data struct {
		MrIID             uint64 `json:"mr_iid"`
		ProjectID         uint64 `json:"project_id"`
		URL               string `json:"url"`
		CurrentCommitHash string `json:"current_commit_hash"`
		NewCommitHash     string `json:"new_commit_hash"`
	} `json:"data"`
}
