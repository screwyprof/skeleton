//go:build acceptance
// +build acceptance

package api_test

import (
	"github.com/brianvoe/gofakeit/v4"

	"github.com/screwyprof/skeleton/internal/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/delivery/rest/resp"
)

func (s *TestSuite) TestIssueCertificate() {
	r := s.givenCorrectRequest()
	got := s.whenUserIssuesACertificate(r)
	s.thenCertificateIsCreated(got)
}

func (s *TestSuite) givenCorrectRequest() req.IssueCertificate {
	return req.IssueCertificate{
		CertificateID: gofakeit.UUID(),
		ArtistID:      gofakeit.UUID(),
		Title:         gofakeit.Sentence(5),
		ArtworkType:   "painting",
	}
}

func (s *TestSuite) whenUserIssuesACertificate(r req.IssueCertificate) *resp.ViewCertificate {
	got, err := s.restClient.IssueCertificate(r)
	s.Require().NoError(err)
	return got
}

func (s *TestSuite) thenCertificateIsCreated(got *resp.ViewCertificate) {
	want := &resp.ViewCertificate{
		CertificateID: got.CertificateID,
		Title:         got.Title,
		ArtworkType:   got.ArtworkType,
		ArtistName:    "Shepard Fairey",
	}
	s.Assert().Equal(want, got)
}
