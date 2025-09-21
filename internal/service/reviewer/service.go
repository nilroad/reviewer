package reviewer

import (
	"club/internal/domain"
	"club/internal/event"
	"club/pkg/cursor"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"

	log "git.oceantim.com/backend/packages/golang/go-logger"
	"github.com/go-git/go-git/v6"
)

func (r *Service) Prepare(ctx context.Context, webhook domain.MergeRequestWebhook) error {
	msg, err := json.Marshal(event.Review{Data: struct {
		MrIID             uint64 `json:"mr_iid"`
		ProjectID         uint64 `json:"project_id"`
		URL               string `json:"url"`
		CurrentCommitHash string `json:"current_commit_hash"`
		NewCommitHash     string `json:"new_commit_hash"`
	}{MrIID: webhook.MergeRequestIID, ProjectID: webhook.ProjectID, CurrentCommitHash: webhook.CurrentCommitHash, NewCommitHash: webhook.NewCommitHash, URL: webhook.URL}})
	if err != nil {
		return err
	}

	err = r.reviewQ.Publish(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}

func (r *Service) Review(ctx context.Context, webhook domain.MergeRequestWebhook) error {
	pid := strconv.FormatUint(webhook.ProjectID, 10)
	mrIID := strconv.FormatUint(webhook.MergeRequestIID, 10)

	tmpDir, err := os.MkdirTemp("", "review-pid"+pid+"-mr_iid"+mrIID)
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			r.logger.Error("failed to remove tmp directory", log.J{
				"tmp_dir": tmpDir,
				"error":   err.Error(),
			})
			return
		}
	}()

	_, err = git.PlainCloneContext(ctx, tmpDir, &git.CloneOptions{URL: webhook.URL, Progress: os.Stdout})
	if err != nil {
		return err
	}

	var prompt strings.Builder

	exampleResult := domain.ReviewResult{
		Score: 10,
		Reviews: []domain.Review{
			{
				File:    "file/path",
				Line:    10,
				Comment: "comment_data",
			},
		},
	}

	exampleJson, err := json.Marshal(exampleResult)
	if err != nil {
		return err
	}

	promptValues := domain.PromptValues{
		ExampleJson:       string(exampleJson),
		SourceDir:         tmpDir,
		NewCommitHash:     webhook.NewCommitHash,
		CurrentCommitHash: webhook.CurrentCommitHash,
	}
	err = r.promptTmpl.Execute(&prompt, promptValues)
	if err != nil {
		return err
	}

	err = cursor.New().OutputFormat(cursor.OutFormatText).Interactive().Prompt(prompt.String()).Execute(ctx, &cursor.ExecConfig{
		Dir:    tmpDir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return err
	}

	reviewFile, err := os.Open(tmpDir + "/review.json")
	if err != nil {
		return err
	}

	reviewBytes, err := io.ReadAll(reviewFile)
	if err != nil {
		return err
	}

	var reviewResult domain.ReviewResult
	err = json.Unmarshal(reviewBytes, &reviewResult)
	if err != nil {
		return err
	}

	for _, v := range reviewResult.Reviews {
		var comment strings.Builder
		err := r.reviewResultTmpl.Execute(&comment, v)
		if err != nil {
			return err
		}

		err = r.gitHostClient.CreateMergeRequestComment(ctx, webhook.ProjectID, webhook.MergeRequestIID, comment.String())
		if err != nil {
			r.logger.Error("failed to create review comment", log.J{
				"error":       err.Error(),
				"review_data": v,
			})

			continue
		}
	}

	return nil
}
