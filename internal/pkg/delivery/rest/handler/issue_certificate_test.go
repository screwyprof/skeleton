package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/golang/mock/gomock"
	"github.com/screwyprof/golibs/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/apierr"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/handler"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/mock"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/resp"
	"github.com/screwyprof/skeleton/pkg/cert/command"
	"github.com/screwyprof/skeleton/pkg/cert/query"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
)

func TestCertificateIssuer_Handle(t *testing.T) { // nolint:funlen
	t.Parallel()

	t.Run("certificate creation error occurs, error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rq := &req.IssueCertificate{
			CertificateID: gofakeit.UUID(),
			ArtistID:      gofakeit.UUID(),
			Title:         gofakeit.Sentence(5),
			ArtworkType:   "painting",
		}

		c := command.IssueCertificate{
			ID:          rq.CertificateID,
			ArtistID:    rq.ArtistID,
			Title:       rq.Title,
			ArtworkType: rq.ArtworkType,
		}

		commandHandler := mock.NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().
			Handle(gomock.Any(), c).
			Return(errSomeBadThingHappened)

		h := handler.NewCertificateIssuer(commandHandler, nil)

		// act
		_, err := h.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, errSomeBadThingHappened)
	})

	t.Run("certificate exists, error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rq := &req.IssueCertificate{
			CertificateID: gofakeit.UUID(),
			ArtistID:      gofakeit.UUID(),
			Title:         gofakeit.Sentence(5),
			ArtworkType:   "painting",
		}

		c := command.IssueCertificate{
			ID:          rq.CertificateID,
			ArtistID:    rq.ArtistID,
			Title:       rq.Title,
			ArtworkType: rq.ArtworkType,
		}

		commandHandler := mock.NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().
			Handle(gomock.Any(), c).
			Return(storage.ErrCertificateAlreadyExists)

		h := handler.NewCertificateIssuer(commandHandler, nil)

		// act
		_, err := h.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, storage.ErrCertificateAlreadyExists)
	})

	t.Run("certificate not found, error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()

		rq := &req.IssueCertificate{
			CertificateID: certificateID,
			ArtistID:      gofakeit.UUID(),
			Title:         gofakeit.Sentence(5),
			ArtworkType:   "painting",
		}

		c := command.IssueCertificate{
			ID:          rq.CertificateID,
			ArtistID:    rq.ArtistID,
			Title:       rq.Title,
			ArtworkType: rq.ArtworkType,
		}

		commandHandler := mock.NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().
			Handle(gomock.Any(), c)

		queryRunner := mock.NewMockQueryRunner(ctrl)
		queryRunner.EXPECT().
			RunQuery(gomock.Any(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}).
			Return(storage.ErrCertificateNotFound)

		h := handler.NewCertificateIssuer(commandHandler, queryRunner)

		// act
		_, err := h.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, storage.ErrCertificateNotFound)
	})

	t.Run("an error occurred when getting a certificate, internal system error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()

		rq := &req.IssueCertificate{
			CertificateID: certificateID,
			ArtistID:      gofakeit.UUID(),
			Title:         gofakeit.Sentence(5),
			ArtworkType:   "painting",
		}

		c := command.IssueCertificate{
			ID:          rq.CertificateID,
			ArtistID:    rq.ArtistID,
			Title:       rq.Title,
			ArtworkType: rq.ArtworkType,
		}

		commandHandler := mock.NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().
			Handle(gomock.Any(), c)

		queryRunner := mock.NewMockQueryRunner(ctrl)
		queryRunner.EXPECT().
			RunQuery(gomock.Any(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}).
			Return(errSomeBadThingHappened)

		h := handler.NewCertificateIssuer(commandHandler, queryRunner)

		// act
		_, err := h.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, errSomeBadThingHappened)
	})

	t.Run("certificate issued successfully, valid response returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()

		want := &resp.IssueCertificate{
			CertificateID: certificateID,
			Title:         gofakeit.Sentence(5),
			ArtistName:    gofakeit.Name(),
			ArtworkType:   "painting",
		}

		rq := &req.IssueCertificate{
			CertificateID: want.CertificateID,
			ArtistID:      gofakeit.UUID(),
			Title:         want.Title,
			ArtworkType:   want.ArtworkType,
		}

		c := command.IssueCertificate{
			ID:          rq.CertificateID,
			ArtistID:    rq.ArtistID,
			Title:       rq.Title,
			ArtworkType: rq.ArtworkType,
		}

		rep := report.Certificate{
			ArtistName:  want.ArtistName,
			ArtworkType: want.ArtworkType,
			Title:       want.Title,
		}

		commandHandler := mock.NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().
			Handle(gomock.Any(), c)

		queryRunner := mock.NewMockQueryRunner(ctrl)
		queryRunner.EXPECT().
			RunQuery(gomock.Any(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}).
			SetArg(2, rep)

		h := handler.NewCertificateIssuer(commandHandler, queryRunner)

		// act
		res, err := h.Handle(context.Background(), rq)

		// assert
		require.NoError(t, err)
		assert.Equals(t, want, res)
	})
}

func assertCause(t *testing.T, got, want error) {
	t.Helper()

	var apiErr *apierr.APIError
	ok := errors.As(got, &apiErr)

	require.True(t, ok)
	assert.Equals(t, apiErr.Cause(), want)
}
