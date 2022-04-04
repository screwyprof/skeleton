package apierr_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/skeleton/internal/delivery/rest/apierr"
)

var errSomeBadThingHappened = errors.New("some error")

func TestAPIError(t *testing.T) {
	t.Parallel()

	err := &apierr.APIError{
		ErrStatus: http.StatusInternalServerError,
		ErrCode:   apierr.InternalServerError,
		ErrMsg:    "Internal Server Error",
		ErrExtra: map[string]interface{}{
			"key": "value",
		},
		ErrCause: errSomeBadThingHappened,
	}

	assert.Equal(t, err.ErrStatus, err.Status())
	assert.Equal(t, err.ErrCode, err.Code())
	assert.Equal(t, err.ErrMsg, err.Message())
	assert.Equal(t, err.ErrExtra, err.Extra())
	assert.Equal(t, err.ErrCause, err.Cause())
	assert.Equal(t, fmt.Sprintf("[%d] %s", err.ErrCode, err.ErrMsg), err.Error())
}

func TestWrap(t *testing.T) {
	t.Parallel()

	want := &apierr.APIError{
		ErrStatus: http.StatusBadRequest,
		ErrCode:   apierr.BadRequest,
		ErrMsg:    "Bad Request",
		ErrExtra:  nil,
		ErrCause:  errSomeBadThingHappened,
	}

	got := apierr.Wrap(want.ErrCause, want.ErrCode, want.ErrExtra)

	assert.Equal(t, want, got)
}
