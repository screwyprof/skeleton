package cmdhandler

import "context"

// CommandHandler runs a command and returns the corresponding result.
type CommandHandler interface {
	Handle(ctx context.Context, command interface{}) error
}

// CommandHandlerFn defines a command handler signature.
type CommandHandlerFn func(ctx context.Context, command interface{}) error
