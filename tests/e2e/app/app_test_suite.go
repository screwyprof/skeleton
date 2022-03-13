//go:build acceptance
// +build acceptance

package app

import (
	"github.com/brianvoe/gofakeit/v4"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/screwyprof/skeleton/internal/pkg/app"
	"github.com/screwyprof/skeleton/internal/pkg/app/modcfg"
)

// TestSuite defines a test suite for the application.
type TestSuite struct {
	suite.Suite
	App *fxtest.App
	Cfg *modcfg.Spec
}

// SetupSuite runs once on suit initialization.
func (s *TestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// seed random generator
	gofakeit.Seed(0)

	s.App = fxtest.New(s.T(),
		app.Module,
		fx.Populate(&s.Cfg),
	)
	s.App.RequireStart()
}

func (s *TestSuite) TearDownSuite() {
	s.App.RequireStop()
}
