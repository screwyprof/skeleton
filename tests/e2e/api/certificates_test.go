//go:build acceptance
// +build acceptance

package api_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/restclient"
	"github.com/screwyprof/skeleton/tests/e2e/app"
)

// TestSuiteCertificates initialize test suit.
func TestSuiteCertificates(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// TestSuite tests the certificates API.
type TestSuite struct {
	app.TestSuite
	restClient *restclient.RESTClient
}

// SetupSuite runs once on suit initialization.
func (s *TestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.restClient = restclient.New("http://localhost:" + s.Cfg.HTTPPort)
}
