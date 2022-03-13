//go:build acceptance
// +build acceptance

package api_test

import (
	"github.com/brianvoe/gofakeit/v4"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/resp"
)

func (s *TestSuite) TestViewCertificate() {
	want := s.givenCertificateExists()
	got := s.whenUserViewsACertificate(want.CertificateID)
	s.thenCertificateIsDisplayed(got)
}

func (s *TestSuite) whenUserViewsACertificate(certificateID string) *resp.ViewCertificate {
	got, err := s.restClient.ViewCertificate(certificateID)
	s.Require().NoError(err)
	return got
}

func (s *TestSuite) givenCertificateExists() *resp.ViewCertificate {
	got, err := s.restClient.IssueCertificate(req.IssueCertificate{
		CertificateID: gofakeit.UUID(),
		ArtistID:      gofakeit.UUID(),
		Title:         gofakeit.Sentence(5),
		ArtworkType:   "painting",
	})
	s.Require().NoError(err)
	return got
}

func (s *TestSuite) thenCertificateIsDisplayed(got *resp.ViewCertificate) {
	want := &resp.ViewCertificate{
		CertificateID: got.CertificateID,
		Title:         got.Title,
		ArtworkType:   got.ArtworkType,
		ArtistName:    "Shepard Fairey",
	}
	s.Assert().Equal(want, got)
}
