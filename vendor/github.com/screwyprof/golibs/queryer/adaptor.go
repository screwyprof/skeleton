package queryer

import (
	"context"
	"errors"
	"reflect"
)

// Guard errors
var (
	ErrInvalidQueryRunnerSignature = errors.New("queryRunner must have 3 input params")
	ErrQueryRunnerIsNotAFunction   = errors.New("queryRunner is not a function")
	ErrFirstArgHasInvalidType      = errors.New("first input argument must have context.Context type")
	ErrSecondArgHasInvalidType     = errors.New("second input argument must be a struct")
	ErrThirdArgHasInvalidType      = errors.New("third input argument must be a pointer to a struct")
)

// Adapt transforms a concrete query runner into a generic one.
// A concrete runner function should have 3 arguments:
// - ctx context.Context,
// - q - a query struct,
// - r - a pointer to a report struct.
//
// The returned param must have error type.
// An example signature may look like:
//   func(ctx context.Context, q TestQuery, r *TestReport) error
//
func Adapt(queryRunner interface{}) (QueryRunnerFn, error) {
	queryRunnerType := reflect.TypeOf(queryRunner)
	err := ensureSignatureIsValid(queryRunnerType)
	if err != nil {
		return nil, err
	}

	fn := func(ctx context.Context, q interface{}, r interface{}) error {
		return invokeQueryRunner(queryRunner, ctx, q, r)
	}

	return fn, nil
}

// MustAdapt acts Like Adapt, but panics on error.
func MustAdapt(queryRunner interface{}) QueryRunnerFn {
	queryRunnerFn, err := Adapt(queryRunner)
	if err != nil {
		panic(err)
	}
	return queryRunnerFn
}

func ensureSignatureIsValid(queryRunnerType reflect.Type) error {
	if queryRunnerType.Kind() != reflect.Func {
		return ErrQueryRunnerIsNotAFunction
	}

	if queryRunnerType.NumIn() != 3 {
		return ErrInvalidQueryRunnerSignature
	}

	return ensureParamsHaveValidTypes(queryRunnerType)
}

func ensureParamsHaveValidTypes(queryRunnerType reflect.Type) error {
	if !firstArgIsContext(queryRunnerType) {
		return ErrFirstArgHasInvalidType
	}

	if !secondArgIsStructure(queryRunnerType) {
		return ErrSecondArgHasInvalidType
	}

	if !thirdArgIsAPointerToAStruct(queryRunnerType) {
		return ErrThirdArgHasInvalidType
	}

	return nil
}

func firstArgIsContext(queryRunnerType reflect.Type) bool {
	ctxtInterface := reflect.TypeOf((*context.Context)(nil)).Elem()
	ctx := queryRunnerType.In(0)
	return ctx.Implements(ctxtInterface)
}

func secondArgIsStructure(queryRunnerType reflect.Type) bool {
	return queryRunnerType.In(1).Kind() == reflect.Struct
}

func thirdArgIsAPointerToAStruct(queryRunnerType reflect.Type) bool {
	thirdArg := queryRunnerType.In(2)
	return thirdArg.Kind() == reflect.Ptr && thirdArg.Elem().Kind() == reflect.Struct

}

func invokeQueryRunner(queryRunner interface{}, args ...interface{}) error {
	result := invoke(queryRunner, args...)
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
