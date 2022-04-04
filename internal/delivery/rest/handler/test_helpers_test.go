package handler_test

import "context"

type QueryRunnerSpy struct {
	Fn func(ctx context.Context, query, report interface{}) error
}

func (s QueryRunnerSpy) RunQuery(ctx context.Context, query, report interface{}) error {
	return s.Fn(ctx, query, report)
}

type CommandHandlerSpy struct {
	Fn func(ctx context.Context, command interface{}) error
}

func (h CommandHandlerSpy) Handle(ctx context.Context, command interface{}) error {
	return h.Fn(ctx, command)
}
