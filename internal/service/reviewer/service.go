package reviewer

import (
	"club/internal/domain"
	"club/internal/event"
	"club/pkg/cursor"
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	log "git.oceantim.com/backend/packages/golang/go-logger"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/storage/filesystem"
)

func (r *Service) Prepare(ctx context.Context, mrIID, projectID, currentCommitHash, newCommitHash string) error {
	msg, err := json.Marshal(event.Review{Data: struct {
		MrIID             string `json:"mr_iid"`
		ProjectID         string `json:"project_id"`
		CurrentCommitHash string `json:"current_commit_hash"`
		NewCommitHash     string `json:"new_commit_hash"`
	}{MrIID: mrIID, ProjectID: projectID, CurrentCommitHash: currentCommitHash, NewCommitHash: newCommitHash}})
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

	repo, err := git.PlainCloneContext(ctx, tmpDir, &git.CloneOptions{URL: webhook.URL, Progress: os.Stdout})
	if err != nil {
		return err
	}

	err = r.cursor.OutputFormat(cursor.OutFormatText).Interactive().Execute(ctx)
	if err != nil {
		return err
	}

}
