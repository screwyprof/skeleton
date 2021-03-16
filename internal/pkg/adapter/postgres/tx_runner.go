package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/screwyprof/golibs/gin/middleware/ctxtx"
)

type CtxTxRunner struct {
	db *pg.DB
}

func NewCtxTxRunner(db *pg.DB) *CtxTxRunner {
	return &CtxTxRunner{db: db}
}

func (r *CtxTxRunner) RunInTransaction(ctx context.Context, fx func(tx *pg.Tx) error) error {
	return ctxtx.RunInTransaction(ctx, r.db, fx)
}
