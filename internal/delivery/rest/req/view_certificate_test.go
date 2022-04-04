package req_test

import (
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/gin-gonic/gin"
	"github.com/screwyprof/golibs/assert"

	"github.com/screwyprof/skeleton/internal/delivery/rest/req"
)

func TestViewCertificate_Bind(t *testing.T) {
	t.Parallel()

	t.Run("incorrect request data given, error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		rq := &req.ViewCertificate{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// act
		err := rq.Bind(c)

		// assert
		assert.NotNil(t, err)
	})

	t.Run("valid request data given, no error", func(t *testing.T) {
		t.Parallel()
		// arrange
		rq := &req.ViewCertificate{}

		r := httptest.NewRequest("post", "/", nil)
		r.RequestURI = "/certificate"

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "certificate_id", Value: gofakeit.UUID()}}

		// act
		err := rq.Bind(c)

		// assert
		assert.NoError(t, err)
	})
}
