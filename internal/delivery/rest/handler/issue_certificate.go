package handler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/screwyprof/golibs/cmdhandler"
	"github.com/screwyprof/golibs/gin/middleware/ctxtags"
	"github.com/screwyprof/golibs/queryer"

	"github.com/screwyprof/skeleton/cert/command"
	"github.com/screwyprof/skeleton/cert/query"
	"github.com/screwyprof/skeleton/cert/report"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
	"github.com/screwyprof/skeleton/internal/delivery/rest/apierr"
	"github.com/screwyprof/skeleton/internal/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/delivery/rest/resp"
)

type CertificateIssuer struct {
	queryRunner    queryer.QueryRunner
	commandHandler cmdhandler.CommandHandler
}

func NewCertificateIssuer(handler cmdhandler.CommandHandler, queryRunner queryer.QueryRunner) *CertificateIssuer {
	return &CertificateIssuer{commandHandler: handler, queryRunner: queryRunner}
}

func (h *CertificateIssuer) Handle(ctx context.Context, r *req.IssueCertificate) (*resp.IssueCertificate, error) {
	c := command.IssueCertificate{
		ID:          r.CertificateID,
		ArtistID:    r.ArtistID,
		Title:       r.Title,
		ArtworkType: r.ArtworkType,
	}
	if err := h.commandHandler.Handle(ctx, c); err != nil {
		return nil, h.handleErr(err)
	}

	var rep report.Certificate
	if err := h.queryRunner.RunQuery(ctx, query.ViewCertificate{ID: r.CertificateID}, &rep); err != nil {
		return nil, h.handleErr(err)
	}

	tags := ctxtags.FromContext(ctx)
	tags.Set("certificate_id", c.ID)
	tags.Set("artist_id", c.ArtistID)

	res := &resp.IssueCertificate{
		CertificateID: r.CertificateID,
		Title:         rep.Title,
		ArtworkType:   rep.ArtworkType,
		ArtistName:    rep.ArtistName,
	}

	return res, nil
}

func (h *CertificateIssuer) handleErr(err error) error {
	switch {
	case errors.Is(err, storage.ErrCertificateNotFound):
		return apierr.Wrap(err, apierr.NotFound, nil)
	case errors.Is(err, storage.ErrCertificateAlreadyExists):
		return apierr.Wrap(err, apierr.Conflict, nil)
	default:
		return apierr.Wrap(err, apierr.InternalServerError, nil)
	}
}
