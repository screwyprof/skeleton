//+build integration

package certificate_repository_test

import (
	"context"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/pkg/errors"

	"github.com/screwyprof/skeleton/internal/pkg/adapter/postgres"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
)

func (s *TestSuite) TestCertificateByID_CertificateExists_CertificateReturned() {
	want := s.givenCertificateExists()
	got := s.whenCertificateIsRetrieved(want.ID)
	s.thenValidCertificateIsReturned(want, got)
}

func (s *TestSuite) TestCertificateByID_UnknownErrorOccurs_ErrorReturned() {
	// arrange
	repo := postgres.NewCertificateRepository(&errTXRunner{db: s.DB})

	// act
	_, err := repo.CertificateByID(context.Background(), gofakeit.UUID())

	// assert
	s.Require().NotNil(err)
	s.Assert().Contains(err.Error(), "cannot fetch certificate")
}

func (s *TestSuite) TestCertificateByID_CertificateDoesNotExist_ErrorReturned() {
	_, err := s.repo.CertificateByID(context.Background(), gofakeit.UUID())
	s.Assert().Equal(storage.ErrCertificateNotFound, errors.Cause(err))
}

func (s *TestSuite) givenCertificateExists() report.Certificate {
	c := &issuecert.Certificate{
		ID:          gofakeit.UUID(),
		ArtistID:    gofakeit.UUID(),
		ArtworkType: "painting",
		Title:       gofakeit.Sentence(5),
	}

	// act
	err := s.repo.Store(context.Background(), c)
	s.Assert().NoError(err)

	return report.Certificate{
		ID:          c.ID,
		Title:       c.Title,
		ArtistName:  "Shepard Fairey",
		ArtworkType: c.ArtworkType,
	}
}

func (s *TestSuite) whenCertificateIsRetrieved(ID string) report.Certificate {
	rep, err := s.repo.CertificateByID(context.Background(), ID)
	s.Require().NoError(err)
	return rep
}

func (s *TestSuite) thenValidCertificateIsReturned(want, got report.Certificate) {
	s.Assert().Equal(want, got)

	s.T().Cleanup(s.removeCertificate(want.ID))
}
