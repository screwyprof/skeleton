package fake

import (
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
)

const existingCertificateID = "b8c70bc2-e2f4-400c-87cc-1a8fef30f7fc"

type StorageReporter struct {
	certificate *issuecert.Certificate
}

func (f *StorageReporter) Store(certificate *issuecert.Certificate) error {
	f.certificate = certificate

	return nil
}

func (f *StorageReporter) CertificateByID(certificateID string) (report.Certificate, error) {
	if certificateID == existingCertificateID {
		return report.Certificate{
			ID:          certificateID,
			ArtistName:  "Shepard Fairey",
			ArtworkType: "painting",
			Title:       "Hope",
		}, nil
	}

	return report.Certificate{
		ArtistName:  "Shepard Fairey",
		ArtworkType: f.certificate.ArtworkType,
		Title:       f.certificate.Title,
	}, nil
}
