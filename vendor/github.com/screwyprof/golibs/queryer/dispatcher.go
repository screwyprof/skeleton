package queryer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ErrNotFound indicates that the query runner was not registered for the given query.
var ErrNotFound = errors.New("query runner not found")

// Dispatcher dispatches the given query to a pre-registered query runner.
type Dispatcher struct {
	runners   map[string]QueryRunnerFn
	runnersMu sync.RWMutex
}

// NewDispatcher creates a new instance of Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		runners: make(map[string]QueryRunnerFn),
	}
}

// RunQuery dispatches the given query to a pre-registered query runner and runs the query.
// Implements QueryRunner interface.
func (d *Dispatcher) RunQuery(ctx context.Context, q interface{}, r interface{}) error {
	queryType := reflect.TypeOf(q).Name()
	runner, ok := d.runners[queryType]
	if !ok {
		return fmt.Errorf("queryer for %s query is not found: %w", queryType, ErrNotFound)
	}
	return runner(ctx, q, r)
}

// RegisterQueryRunner registers a query runner.
// TODO: add guards for an empty runner name and a nil runner.
func (d *Dispatcher) RegisterQueryRunner(queryRunnerFnName string, queryRunnerFn QueryRunnerFn) {
	d.runnersMu.Lock()
	defer d.runnersMu.Unlock()

	d.runners[queryRunnerFnName] = queryRunnerFn
}
