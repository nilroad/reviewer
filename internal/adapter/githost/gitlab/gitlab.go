package gitlab

import (
	"context"

	"git.oceantim.com/backend/packages/golang/essential/exception"
	"git.oceantim.com/backend/packages/golang/go-utils"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (r *Gitlab) CreateMergeRequestComment(ctx context.Context, pid uint64, mrIID uint64, body string) error {
	_, res, err := r.client.Discussions.CreateMergeRequestDiscussion(pid, int(mrIID), &gitlab.CreateMergeRequestDiscussionOptions{
		Body: &body,
	})
	if err != nil {
		return exception.NewUnExpectedError(ctx, err)
	}

	if !utils.IsHttpRequestSuccess(res.StatusCode) {
		return exception.NewUnExpectedError(ctx, err)
	}

	return nil
}
