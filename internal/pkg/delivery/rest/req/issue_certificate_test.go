package req_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/gin-gonic/gin"
	"github.com/screwyprof/golibs/assert"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/req"
)

func TestIssueCertificate_Bind(t *testing.T) {
	t.Parallel()

	t.Run("incorrect request data given, error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		rq := &req.IssueCertificate{}

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
		rq := &req.IssueCertificate{
			CertificateID: gofakeit.UUID(),
			ArtistID:      gofakeit.UUID(),
			Title:         gofakeit.Sentence(5),
			ArtworkType:   "painting",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("post", "/", strings.NewReader(assert.JSONStringFor(t, rq)))

		// act
		err := rq.Bind(c)

		// assert
		assert.NoError(t, err)
	})
}
