package viewcert

import (
	"context"
	"fmt"

	"github.com/screwyprof/skeleton/pkg/cert/query"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
)

type CertReporter interface {
	CertificateByID(ctx context.Context, certificateID string) (report.Certificate, error)
}

type CertViewer struct {
	reporter CertReporter
}

func NewQueryer(reporter CertReporter) *CertViewer {
	return &CertViewer{reporter: reporter}
}

func (v CertViewer) ViewCertificate(ctx context.Context, q query.ViewCertificate, r *report.Certificate) error {
	cert, err := v.reporter.CertificateByID(ctx, q.ID)
	if err != nil {
		return fmt.Errorf("certificateID: %s: %w", q.ID, storage.ErrCertificateNotFound)
	}

	r.Title = cert.Title
	r.ArtistName = cert.ArtistName
	r.ArtworkType = cert.ArtworkType

	return nil
}
