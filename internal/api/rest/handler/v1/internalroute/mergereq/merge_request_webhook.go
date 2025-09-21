package mergereq

import (
	"club/internal/api/rest/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Handler) MergeRequestWebhook(c *gin.Context) {
	var body *request.MRWebhook
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = r.validator.Validate(
		c,
		body,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = r.reviewService.Prepare(c, body.ToDomain())
	if err != nil {
		return
	}

	c.Status(http.StatusAccepted)
}
