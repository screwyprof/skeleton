package assert

import (
	"encoding/json"
	"reflect"
	"testing"
)

// True fails the test if the condition is false.
func True(tb testing.TB, condition bool) {
	tb.Helper()
	if !condition {
		tb.Errorf("\033[31mcondition is false\033[39m\n\n")
	}
}

// NoError fails the test if an err is not nil.
func NoError(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("\033[31munexpected error: %v\033[39m\n\n", err)
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Errorf("\033[31m\n\n\texp:\n\t%#+v\n\n\tgot:\n\t%#+v\033[39m", exp, act)
	}
}

// NotNil fails the test if the argument is nil.
func NotNil(tb testing.TB, arg interface{}) {
	tb.Helper()
	if arg == nil {
		tb.Fatal("\033[31mexpected non nil value, but got: <nil>\033[39m")
	}
}

// Panic fails the test if it didn't panic.
func Panic(tb testing.TB, f func()) {
	tb.Helper()
	defer func() {
		tb.Helper()
		if r := recover(); r == nil {
			tb.Errorf("\033[31mpanic is expected\033[39m")
		}
	}()
	f()
}

// JSONStringFor returns json-string for the given object or fails the test.
func JSONStringFor(tb testing.TB, object interface{}) string {
	tb.Helper()
	bytes, err := json.Marshal(object)
	if err != nil {
		tb.Errorf("\033[31mcannot marshal the given object\033[39m")
	}
	return string(bytes)
}
