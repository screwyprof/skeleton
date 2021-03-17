package req_test

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	gofakeit.Seed(0)

	exitVal := m.Run()
	os.Exit(exitVal)
}
