package cmdhandler

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ErrNotFound indicates that the command handler was not registered for the given command.
var ErrNotFound = errors.New("command handler not found")

// Dispatcher dispatches the given command to a pre-registered command handler.
type Dispatcher struct {
	handlers   map[string]CommandHandlerFn
	handlersMu sync.RWMutex
}

// NewDispatcher creates a new instance of Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]CommandHandlerFn),
	}
}

// Handle dispatches the given command to a pre-registered command handler and runs the command.
// Implements CommandHandler interface.
func (d *Dispatcher) Handle(ctx context.Context, c interface{}) error {
	commandType := reflect.TypeOf(c).Name()
	runner, ok := d.handlers[commandType]
	if !ok {
		return fmt.Errorf("cmdhandler for %s command is not found: %w", commandType, ErrNotFound)
	}
	return runner(ctx, c)
}

// RegisterCommandHandler registers a command handler.
// TODO: add guards for an empty runner name and a nil runner.
func (d *Dispatcher) RegisterCommandHandler(commandHandlerFnName string, commandHandlerFn CommandHandlerFn) {
	d.handlersMu.Lock()
	defer d.handlersMu.Unlock()

	d.handlers[commandHandlerFnName] = commandHandlerFn
}
