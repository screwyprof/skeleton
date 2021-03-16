package ctxtx

import (
	"context"

	"github.com/go-pg/pg/v9"
)

func RunInTransaction(ctx context.Context, db *pg.DB, fn func(tx *pg.Tx) error) error {
	tx, isInSession, err := extractTX(ctx, db)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if isInSession {
			return err
		}
		_ = tx.Rollback()
		return err
	}

	if !isInSession {
		return tx.Commit()
	}
	return nil
}

func extractTX(ctx context.Context, db *pg.DB) (*pg.Tx, bool, error) {
	tx := FromContext(ctx)
	if tx != nil {
		return tx, true, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, false, err
	}
	return tx, false, nil
}
