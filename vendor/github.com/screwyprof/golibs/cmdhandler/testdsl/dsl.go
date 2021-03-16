package testdsl

import (
	"context"
	"errors"
	"testing"

	"github.com/screwyprof/golibs/assert"
	"github.com/screwyprof/golibs/cmdhandler"
)

// GivenFn is a test init function.
type GivenFn func() (cmdhandler.CommandHandler, error)

// WhenFn is a command handler function.
type WhenFn func(commandHandler cmdhandler.CommandHandler, err error) error

// ThenFn prepares the Checker.
type ThenFn func(t testing.TB) Checker

// Checker asserts the given results.
type Checker func(err error)

// DispatcherTester defines a query runner tester.
type DispatcherTester func(given GivenFn, when WhenFn, then ThenFn)

// Test runs the test.
//
// Example:
//	commandHandler := &mock.ConcreteCommandHandlerStub{}
//
//	Test(t)(
//		Given("TestCommand", commandHandler.Handle),
//		When(context.Background(), mock.TestCommand{}),
//		ThenOk(),
//	)
//
func Test(t testing.TB) DispatcherTester {
	return func(given GivenFn, when WhenFn, then ThenFn) {
		t.Helper()
		then(t)(when(given()))
	}
}

// Given prepares the given command handler for testing.
func Given(commandType string, concreteCommandHandler interface{}) GivenFn {
	return func() (cmdhandler.CommandHandler, error) {
		commandHandler, err := createCommandHandler(commandType, concreteCommandHandler)
		if err != nil {
			return nil, err
		}
		return commandHandler, nil
	}
}

// When prepares the command handler for the given command.
func When(ctx context.Context, c interface{}) WhenFn {
	return func(commandHandler cmdhandler.CommandHandler, runnerError error) error {
		if runnerError != nil {
			return runnerError
		}
		return commandHandler.Handle(ctx, c)
	}
}

// Then asserts that the expected events are applied.
func ThenOk() ThenFn {
	return func(t testing.TB) Checker {
		return func(err error) {
			t.Helper()
			assert.NoError(t, err)
		}
	}
}

// ThenFailWith asserts that the expected error occurred.
func ThenFailWith(want error) ThenFn {
	return func(t testing.TB) Checker {
		return func(err error) {
			t.Helper()
			assert.NotNil(t, err)
			assert.True(t, errors.Is(err, want))
		}
	}
}

func createCommandHandler(commandHandlerName string, h interface{}) (*cmdhandler.Dispatcher, error) {
	commandHandler, err := cmdhandler.Adapt(h)
	if err != nil {
		return nil, err
	}

	dispatcher := cmdhandler.NewDispatcher()
	dispatcher.RegisterCommandHandler(commandHandlerName, commandHandler)
	return dispatcher, nil
}
