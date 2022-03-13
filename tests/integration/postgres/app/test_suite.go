//go:build integration

package app

import (
	"github.com/brianvoe/gofakeit/v4"
	"github.com/gin-gonic/gin"
	"github.com/go-rel/rel"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// TestSuite defines a test suite for the application.
type TestSuite struct {
	suite.Suite
	app  *fxtest.App
	Repo rel.Repository
}

// SetupSuite runs once on suit initialization.
func (s *TestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// seed random generator
	gofakeit.Seed(0)

	s.app = fxtest.New(s.T(),
		TestModule,
		fx.Populate(&s.Repo),
	)
	s.app.RequireStart()
}

func (s *TestSuite) TearDownSuite() {
	s.app.RequireStop()
}
