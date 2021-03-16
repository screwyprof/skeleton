package ctxtx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
)

func CtxTX(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := db.RunInTransaction(func(tx *pg.Tx) error {
			c = ToContext(c, tx)
			c.Next()
			return nil
		})

		if err != nil {
			_ = c.AbortWithError(c.Writer.Status(), err)
		}
	}
}
