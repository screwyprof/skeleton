//+build integration

package certificate_repository_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/screwyprof/skeleton/internal/pkg/adapter/postgres"
	"github.com/screwyprof/skeleton/internal/pkg/adapter/postgres/model"
	"github.com/screwyprof/skeleton/tests/integration/postgres/app"
)

// TestSuiteCertificateRepository initialize test suit.
func TestSuiteCertificateRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// TestSuite tests the certificate repository.
type TestSuite struct {
	app.TestSuite
	repo *postgres.CertificateRepository
}

// SetupSuite runs once on suit initialization.
func (s *TestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.repo = postgres.NewCertificateRepository(postgres.NewCtxTxRunner(s.DB))
}

func (s *TestSuite) certificateByID(ID string) model.Certificate {
	var row model.Certificate
	err := s.DB.Model(&row).Where("certificate_id = ?", ID).Select()
	s.Require().NoError(err)
	return row
}

func (s *TestSuite) removeCertificate(certificateID string) func() {
	return func() {
		c := model.Certificate{
			CertificateID: certificateID,
		}
		_, err := s.DB.Model(&c).WherePK().Delete()
		s.Require().NoError(err)
	}
}
