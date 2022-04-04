//go:build integration
// +build integration

package certificate_repository_test

import (
	"context"

	"github.com/brianvoe/gofakeit/v4"

	"github.com/screwyprof/skeleton/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
	"github.com/screwyprof/skeleton/internal/adapter/postgres/model"
)

func (s *TestSuite) TestStore_ValidDataGiven_CertificateCreated() {
	c := s.givenCertificate()
	s.whenCertificateIsStored(c)
	s.thenCertificateExists(c)
}

func (s *TestSuite) TestStore_CannotStoreCertificate_ErrorReturned() {
	err := s.repo.Store(context.Background(), &issuecert.Certificate{})
	s.Assert().ErrorIs(err, storage.ErrCannotStoreCertificate)
}

func (s *TestSuite) TestStore_SameCertificateStoredTwice_ErrorReturned() {
	c := s.givenCertificate()
	err := s.whenCertificateIsStoredTwice(c)
	s.thenFailsWithAlreadyExistsError(err)
}

func (s *TestSuite) givenCertificate() *issuecert.Certificate {
	c := &issuecert.Certificate{
		ID:          gofakeit.UUID(),
		ArtistID:    gofakeit.UUID(),
		ArtworkType: "painting",
		Title:       gofakeit.Sentence(5),
	}
	return c
}

func (s *TestSuite) whenCertificateIsStored(c *issuecert.Certificate) {
	err := s.repo.Store(context.Background(), c)
	s.Assert().NoError(err)
}

func (s *TestSuite) whenCertificateIsStoredTwice(c *issuecert.Certificate) error {
	err := s.repo.Store(context.Background(), c)
	s.Assert().NoError(err)

	return s.repo.Store(context.Background(), c)
}

func (s *TestSuite) thenFailsWithAlreadyExistsError(err error) {
	s.Assert().ErrorIs(err, storage.ErrCertificateAlreadyExists)
}

func (s *TestSuite) thenCertificateExists(c *issuecert.Certificate) {
	row := s.certificateByID(c.ID)

	want := model.Certificate{
		CertificateID: c.ID,
		ArtistID:      c.ArtistID,
		Title:         c.Title,
		ArtworkType:   c.ArtworkType,
	}
	s.Assert().Equal(want, row)

	s.T().Cleanup(s.removeCertificate(c.ID))
}
