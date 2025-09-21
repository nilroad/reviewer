package githost

import "context"

type Contract interface {
	CreateMergeRequestComment(ctx context.Context, pid uint64, mrIID uint64, body string) error
}
