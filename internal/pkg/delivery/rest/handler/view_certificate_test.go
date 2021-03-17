package handler_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/golang/mock/gomock"
	"github.com/screwyprof/golibs/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/handler"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/mock"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/resp"
	"github.com/screwyprof/skeleton/pkg/cert/query"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
)

func TestCertificateViewer_Handle(t *testing.T) { // nolint:funlen
	t.Parallel()

	t.Run("certificate does not exist, error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()

		queryRunner := mock.NewMockQueryRunner(ctrl)
		queryRunner.EXPECT().
			RunQuery(context.Background(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}).
			Return(storage.ErrCertificateNotFound)

		h := handler.NewCertificateViewer(queryRunner)

		// act
		_, err := h.Handle(context.Background(), &req.ViewCertificate{CertificateID: certificateID})

		// assert
		assertCause(t, err, storage.ErrCertificateNotFound)
	})

	t.Run("unknown error occurred, internal system error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()

		queryRunner := mock.NewMockQueryRunner(ctrl)
		queryRunner.EXPECT().
			RunQuery(context.Background(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}).
			Return(errSomeBadThingHappened)

		h := handler.NewCertificateViewer(queryRunner)

		// act
		_, err := h.Handle(context.Background(), &req.ViewCertificate{CertificateID: certificateID})

		// assert
		assertCause(t, err, errSomeBadThingHappened)
	})

	t.Run("certificate exists, valid response returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()

		want := &resp.ViewCertificate{
			CertificateID: certificateID,
			Title:         gofakeit.Sentence(5),
			ArtistName:    gofakeit.Name(),
			ArtworkType:   "painting",
		}

		certificate := report.Certificate{
			ArtistName:  want.ArtistName,
			ArtworkType: want.ArtworkType,
			Title:       want.Title,
		}

		queryRunner := mock.NewMockQueryRunner(ctrl)
		queryRunner.EXPECT().
			RunQuery(context.Background(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}).
			SetArg(2, certificate)

		h := handler.NewCertificateViewer(queryRunner)

		// act
		res, err := h.Handle(context.Background(), &req.ViewCertificate{CertificateID: certificateID})

		// assert
		require.NoError(t, err)
		assert.Equals(t, want, res)
	})
}
