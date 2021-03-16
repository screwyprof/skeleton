package adaptor

import (
	"context"
	"reflect"
)

type HandlerFn func(context.Context, interface{}) (interface{}, error)

func Adapt(h interface{}) (HandlerFn, error) {
	fn := func(ctx context.Context, r interface{}) (interface{}, error) {
		return invokeHandler(h, ctx, r)
	}

	return fn, nil
}

// MustAdapt acts Like Adapt, but panics on error.
func MustAdapt(h interface{}) HandlerFn {
	fn, err := Adapt(h)
	if err != nil {
		panic(err)
	}
	return fn
}

func invokeHandler(h interface{}, args ...interface{}) (interface{}, error) {
	result := invoke(h, args...)
	resErr := result[1].Interface()
	if resErr != nil {
		return nil, resErr.(error)
	}
	return result[0].Interface(), nil
}

func invoke(fn interface{}, args ...interface{}) []reflect.Value {
	v := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	return v.Call(in)
}
