package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"

	"github.com/screwyprof/skeleton/internal/pkg/adapter/postgres/model"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
)

type TxRunner interface {
	RunInTransaction(ctx context.Context, fx func(tx *pg.Tx) error) error
}

type CertificateRepository struct {
	txRunner TxRunner
}

func NewCertificateRepository(txRunner TxRunner) *CertificateRepository {
	return &CertificateRepository{txRunner: txRunner}
}

func (r *CertificateRepository) CertificateByID(ctx context.Context, certificateID string) (report.Certificate, error) {
	var row model.Certificate

	err := r.txRunner.RunInTransaction(ctx, func(tx *pg.Tx) error {
		var err error
		row, err = r.fetchCertificate(tx, certificateID)

		return err
	})
	if err != nil {
		return report.Certificate{}, err
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
	return r.txRunner.RunInTransaction(ctx, func(tx *pg.Tx) error {
		return r.createCertificate(tx, model.Certificate{
			CertificateID: c.ID,
			ArtistID:      c.ArtistID,
			Title:         c.Title,
			ArtworkType:   c.ArtworkType,
		})
	})
}

func (r *CertificateRepository) createCertificate(tx *pg.Tx, row model.Certificate) error {
	_, err := tx.Model(&row).Insert()
	if err != nil {
		if isErrUniqueConstraintViolation(err) {
			return errors.Wrap(storage.ErrCertificateAlreadyExists, "id="+row.CertificateID)
		}

		return errors.Wrap(storage.ErrCannotStoreCertificate, "id="+row.CertificateID)
	}

	return nil
}

func (r *CertificateRepository) fetchCertificate(tx *pg.Tx, certificateID string) (model.Certificate, error) {
	var row model.Certificate

	err := tx.Model(&row).Where("certificate_id = ?", certificateID).Select()
	if errors.Is(err, pg.ErrNoRows) {
		return model.Certificate{}, errors.Wrap(
			storage.ErrCertificateNotFound, "certificate "+certificateID+" is not found")
	}

	if err != nil {
		return model.Certificate{}, errors.Wrap(err, "cannot fetch certificate "+certificateID)
	}

	return row, nil
}

func isErrUniqueConstraintViolation(err error) bool {
	var pgErr pg.Error
	if ok := errors.As(err, &pgErr); ok {
		if pgErr.Field('C') == "23505" {
			return true
		}
	}

	return false
}
