package gitlab

import "gitlab.com/gitlab-org/api/client-go"

type Gitlab struct {
	client *gitlab.Client
}

func New(client *gitlab.Client) *Gitlab {
	return &Gitlab{
		client: client,
	}
}
