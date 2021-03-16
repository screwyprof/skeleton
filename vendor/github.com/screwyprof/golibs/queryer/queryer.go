package queryer

import "context"

// QueryRunner runs a query and returns the corresponding result setting report value by ref.
type QueryRunner interface {
	RunQuery(ctx context.Context, query interface{}, report interface{}) error
}

// QueryRunnerFn defines a query runner signature.
type QueryRunnerFn func(ctx context.Context, query interface{}, report interface{}) error
