package handler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/screwyprof/golibs/queryer"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/apierr"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/resp"
	"github.com/screwyprof/skeleton/pkg/cert/query"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
)

type CertificateViewer struct {
	queryRunner queryer.QueryRunner
}

func NewCertificateViewer(queryRunner queryer.QueryRunner) *CertificateViewer {
	return &CertificateViewer{queryRunner: queryRunner}
}

func (h *CertificateViewer) Handle(ctx context.Context, r *req.ViewCertificate) (*resp.ViewCertificate, error) {
	var rep report.Certificate

	q := query.ViewCertificate{ID: r.CertificateID}
	if err := h.queryRunner.RunQuery(ctx, q, &rep); err != nil {
		return nil, h.handleErr(err)
	}

	res := &resp.ViewCertificate{
		CertificateID: r.CertificateID,
		Title:         rep.Title,
		ArtworkType:   rep.ArtworkType,
		ArtistName:    rep.ArtistName,
	}

	return res, nil
}

func (h *CertificateViewer) handleErr(err error) error {
	switch {
	case errors.Is(err, storage.ErrCertificateNotFound):
		return apierr.Wrap(err, apierr.NotFound, nil)
	default:
		return apierr.Wrap(err, apierr.InternalServerError, nil)
	}
}
