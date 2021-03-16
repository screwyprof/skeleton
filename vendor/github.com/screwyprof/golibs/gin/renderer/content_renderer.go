package renderer

import "github.com/gin-gonic/gin"

type ContentRenderer struct {
	Status  int
	Content interface{}
}

func (r *ContentRenderer) Render(c *gin.Context) {
	c.JSON(r.Status, r.Content)
}
