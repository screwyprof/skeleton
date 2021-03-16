package cmdhandler

import (
	"context"
	"errors"
	"reflect"
)

// Guard errors
var (
	ErrInvalidCommandHandlerSignature = errors.New("command handler must have 3 input params")
	ErrCommandHandlerIsNotAFunction   = errors.New("command handler is not a function")
	ErrFirstArgHasInvalidType         = errors.New("first input argument must have context.Context type")
	ErrSecondArgHasInvalidType        = errors.New("second input argument must be a struct")
)

// Adapt transforms a concrete command handler into a generic one.
// A concrete runner function should have 3 arguments:
// - ctx context.Context,
// - c - a command struct,
//
// The returned param must have error type.
// An example signature may look like:
//   func(ctx context.Context, c TestCommand) error
//
func Adapt(commandHandler interface{}) (CommandHandlerFn, error) {
	commandHandlerType := reflect.TypeOf(commandHandler)
	err := ensureSignatureIsValid(commandHandlerType)
	if err != nil {
		return nil, err
	}

	fn := func(ctx context.Context, q interface{}) error {
		return invokeCommandHandler(commandHandler, ctx, q)
	}

	return fn, nil
}

// MustAdapt acts Like Adapt, but panics on error.
func MustAdapt(commandHandler interface{}) CommandHandlerFn {
	h, err := Adapt(commandHandler)
	if err != nil {
		panic(err)
	}
	return h
}

func ensureSignatureIsValid(commandHandlerType reflect.Type) error {
	if commandHandlerType.Kind() != reflect.Func {
		return ErrCommandHandlerIsNotAFunction
	}

	if commandHandlerType.NumIn() != 2 {
		return ErrInvalidCommandHandlerSignature
	}

	return ensureParamsHaveValidTypes(commandHandlerType)
}

func ensureParamsHaveValidTypes(commandHandlerType reflect.Type) error {
	if !firstArgIsContext(commandHandlerType) {
		return ErrFirstArgHasInvalidType
	}

	if !secondArgIsStructure(commandHandlerType) {
		return ErrSecondArgHasInvalidType
	}

	return nil
}

func firstArgIsContext(commandHandlerType reflect.Type) bool {
	ctxtInterface := reflect.TypeOf((*context.Context)(nil)).Elem()
	ctx := commandHandlerType.In(0)
	return ctx.Implements(ctxtInterface)
}

func secondArgIsStructure(commandHandlerType reflect.Type) bool {
	return commandHandlerType.In(1).Kind() == reflect.Struct
}

func invokeCommandHandler(commandHandler interface{}, args ...interface{}) error {
	result := invoke(commandHandler, args...)
	resErr := result[0].Interface()
	if resErr != nil {
		return resErr.(error)
	}
	return nil
}

func invoke(fn interface{}, args ...interface{}) []reflect.Value {
	v := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	return v.Call(in)
}
