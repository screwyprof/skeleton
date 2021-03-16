package ctxtx

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
)

var ctxKey = "ctxtx"

func ToContext(c *gin.Context, tx *pg.Tx) *gin.Context {
	c.Set(ctxKey, tx)
	return c
}

func FromContext(ctx context.Context) *pg.Tx {
	if tx, ok := ctx.Value(ctxKey).(*pg.Tx); ok {
		return tx
	}
	return nil
}
