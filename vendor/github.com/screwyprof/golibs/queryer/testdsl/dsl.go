package testdsl

import (
	"context"
	"errors"
	"testing"

	"github.com/screwyprof/golibs/assert"
	"github.com/screwyprof/golibs/queryer"
)

// GivenFn is a test init function.
type GivenFn func() (queryer.QueryRunner, error)

// WhenFn is a command handler function.
type WhenFn func(queryRunner queryer.QueryRunner, err error) (interface{}, error)

// ThenFn prepares the Checker.
type ThenFn func(t testing.TB) Checker

// Checker asserts the given results.
type Checker func(report interface{}, err error)

// QueryerTester defines a query runner tester.
type QueryerTester func(given GivenFn, when WhenFn, then ThenFn)

// Test runs the test.
//
// Example:
//	want := &mock.TestReport{Value: 123}
//	queryRunner := &mock.ConcreteQueryerStub{Rep: 123}
//
//	Test(t)(
//		Given("TestQuery", queryRunner.Run),
//		When(context.Background(), mock.TestQuery{}, &mock.TestReport{}),
//		Then(want),
//	)
//
func Test(t testing.TB) QueryerTester {
	return func(given GivenFn, when WhenFn, then ThenFn) {
		t.Helper()
		then(t)(when(given()))
	}
}

// Given prepares the given query runner for testing.
func Given(queryType string, concreteQueryRunner interface{}) GivenFn {
	return func() (queryer.QueryRunner, error) {
		queryRunner, err := createQueryRunner(queryType, concreteQueryRunner)
		if err != nil {
			return nil, err
		}
		return queryRunner, nil
	}
}

// When prepares the command handler for the given command.
func When(ctx context.Context, q interface{}, r interface{}) WhenFn {
	return func(queryRunner queryer.QueryRunner, runnerError error) (interface{}, error) {
		if runnerError != nil {
			return nil, runnerError
		}
		return r, queryRunner.RunQuery(ctx, q, r)
	}
}

// Then asserts that the expected events are applied.
func Then(want interface{}) ThenFn {
	return func(t testing.TB) Checker {
		return func(got interface{}, err error) {
			t.Helper()
			assert.NoError(t, err)
			assert.Equals(t, want, got)
		}
	}
}

// ThenFailWith asserts that the expected error occurred.
func ThenFailWith(want error) ThenFn {
	return func(t testing.TB) Checker {
		return func(report interface{}, err error) {
			t.Helper()
			assert.NotNil(t, err)
			assert.True(t, errors.Is(err, want))
		}
	}
}

func createQueryRunner(queryRunnerName string, concreteQueryRunner interface{}) (*queryer.Dispatcher, error) {
	queryRunner, err := queryer.Adapt(concreteQueryRunner)
	if err != nil {
		return nil, err
	}

	dispatcher := queryer.NewDispatcher()
	dispatcher.RegisterQueryRunner(queryRunnerName, queryRunner)
	return dispatcher, nil
}
