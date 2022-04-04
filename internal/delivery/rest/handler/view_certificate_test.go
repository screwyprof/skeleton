package handler_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/screwyprof/golibs/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/skeleton/cert/report"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
	"github.com/screwyprof/skeleton/internal/delivery/rest/handler"
	"github.com/screwyprof/skeleton/internal/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/delivery/rest/resp"
)

func TestCertificateViewer(t *testing.T) {
	t.Parallel()

	t.Run("certificate does not exist, error returned", func(t *testing.T) {
		t.Parallel()

		// arrange
		certificateID := gofakeit.UUID()

		queryRunner := QueryRunnerSpy{
			Fn: func(ctx context.Context, query, report interface{}) error {
				return storage.ErrCertificateNotFound
			},
		}

		sut := handler.NewCertificateViewer(queryRunner)

		// act
		_, err := sut.Handle(context.Background(), &req.ViewCertificate{CertificateID: certificateID})

		// assert
		assertCause(t, err, storage.ErrCertificateNotFound)
	})

	t.Run("unknown error occurred, internal system error returned", func(t *testing.T) {
		t.Parallel()

		// arrange
		certificateID := gofakeit.UUID()

		queryRunner := QueryRunnerSpy{
			Fn: func(ctx context.Context, query, report interface{}) error {
				return errSomeBadThingHappened
			},
		}

		sut := handler.NewCertificateViewer(queryRunner)

		// act
		_, err := sut.Handle(context.Background(), &req.ViewCertificate{CertificateID: certificateID})

		// assert
		assertCause(t, err, errSomeBadThingHappened)
	})

	t.Run("certificate exists, valid response returned", func(t *testing.T) {
		t.Parallel()

		// arrange
		var want *resp.ViewCertificate
		gofakeit.Struct(&want)

		certificate := report.Certificate{
			ID:          want.CertificateID,
			ArtistName:  want.ArtistName,
			ArtworkType: want.ArtworkType,
			Title:       want.Title,
		}

		queryRunner := QueryRunnerSpy{
			Fn: func(ctx context.Context, query, rep interface{}) error {
				r, ok := rep.(*report.Certificate)
				assert.True(t, ok)

				*r = certificate

				return nil
			},
		}

		sut := handler.NewCertificateViewer(queryRunner)

		// act
		res, err := sut.Handle(context.Background(), &req.ViewCertificate{CertificateID: want.CertificateID})

		// assert
		require.NoError(t, err)
		assert.Equals(t, want, res)
	})
}
