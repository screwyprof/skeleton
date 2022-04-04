package postgres

import (
	"context"
	"errors"

	"github.com/ansel1/merry/v2"
	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"

	"github.com/screwyprof/skeleton/internal/adapter/postgres/model"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
)

type CertificateRepository struct {
	repo rel.Repository
}

func NewCertificateRepository(repo rel.Repository) *CertificateRepository {
	return &CertificateRepository{repo: repo}
}

func (r *CertificateRepository) CertificateByID(ctx context.Context, certificateID string) (report.Certificate, error) {
	var row model.Certificate

	if err := r.repo.Find(ctx, &row, where.Eq("certificate_id", certificateID)); err != nil {
		if errors.Is(err, rel.NotFoundError{}) {
			return report.Certificate{},
				merry.Wrap(storage.ErrCertificateNotFound, merry.WithValue("certificate_id", certificateID))
		}

		return report.Certificate{}, merry.Wrap(err)
	}

	c := report.Certificate{
		ID:          row.CertificateID,
		Title:       row.Title,
		ArtistName:  "Shepard Fairey",
		ArtworkType: row.ArtworkType,
	}

	return c, nil
}

func (r *CertificateRepository) Store(ctx context.Context, c *issuecert.Certificate) error {
	row := model.Certificate{
		CertificateID: c.ID,
		ArtistID:      c.ArtistID,
		Title:         c.Title,
		ArtworkType:   c.ArtworkType,
	}

	err := r.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := r.createCertificate(ctx, row); err != nil {
			return merry.Wrap(err)
		}

		return nil
	})
	if err != nil {
		return merry.Wrap(storage.ErrCannotStoreCertificate, merry.WithCause(err))
	}

	return nil
}

func (r *CertificateRepository) createCertificate(ctx context.Context, row model.Certificate) error {
	if err := r.repo.Insert(ctx, &row); err != nil {
		if errors.Is(err, rel.ErrUniqueConstraint) {
			return merry.Wrap(storage.ErrCertificateAlreadyExists,
				merry.WithValue("certificate_id", row.CertificateID))
		}

		return merry.Wrap(err)
	}

	return nil
}
