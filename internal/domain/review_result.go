package domain

type Review struct {
	File    string `json:"file"`
	Line    uint64 `json:"line"`
	Comment string `json:"comment"`
}
type ReviewResult struct {
	Score   float64  `json:"score"`
	Reviews []Review `json:"reviews"`
}
