//+build integration

package certificate_repository_test

import (
	"context"

	"github.com/go-pg/pg/v9"
)

type errTXRunner struct {
	db *pg.DB
}

func (r *errTXRunner) RunInTransaction(_ context.Context, fx func(tx *pg.Tx) error) error {
	return r.db.RunInTransaction(func(tx *pg.Tx) error {
		_ = tx.Commit() // break the transaction
		return fx(tx)
	})
}
