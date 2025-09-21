package event

type Review struct {
	Data struct {
		MrIID             string `json:"mr_iid"`
		ProjectID         string `json:"project_id"`
		CurrentCommitHash string `json:"current_commit_hash"`
		NewCommitHash     string `json:"new_commit_hash"`
	} `json:"data"`
}
