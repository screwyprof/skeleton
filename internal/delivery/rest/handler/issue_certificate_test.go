package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/screwyprof/golibs/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/skeleton/cert/command"
	"github.com/screwyprof/skeleton/cert/report"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
	"github.com/screwyprof/skeleton/internal/delivery/rest/apierr"
	"github.com/screwyprof/skeleton/internal/delivery/rest/handler"
	"github.com/screwyprof/skeleton/internal/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/delivery/rest/resp"
)

func TestCertificateIssuer(t *testing.T) {
	t.Parallel()

	t.Run("certificate creation error occurs, error returned", func(t *testing.T) {
		t.Parallel()

		// arrange
		var rq *req.IssueCertificate
		gofakeit.Struct(&rq)

		commandHandler := CommandHandlerSpy{
			Fn: func(ctx context.Context, command interface{}) error {
				return errSomeBadThingHappened
			},
		}

		sut := handler.NewCertificateIssuer(commandHandler, nil)

		// act
		_, err := sut.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, errSomeBadThingHappened)
	})

	t.Run("certificate exists, error returned", func(t *testing.T) {
		t.Parallel()

		// arrange
		var rq *req.IssueCertificate
		gofakeit.Struct(&rq)

		commandHandler := CommandHandlerSpy{
			Fn: func(ctx context.Context, command interface{}) error {
				return storage.ErrCertificateAlreadyExists
			},
		}

		sut := handler.NewCertificateIssuer(commandHandler, nil)

		// act
		_, err := sut.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, storage.ErrCertificateAlreadyExists)
	})

	t.Run("certificate not found, error returned", func(t *testing.T) {
		t.Parallel()

		// arrange
		certificateID := gofakeit.UUID()

		rq := &req.IssueCertificate{
			CertificateID: certificateID,
			ArtistID:      gofakeit.UUID(),
			Title:         gofakeit.Sentence(5),
			ArtworkType:   "painting",
		}

		commandHandler := CommandHandlerSpy{
			Fn: func(ctx context.Context, command interface{}) error {
				return nil
			},
		}

		queryRunner := QueryRunnerSpy{
			Fn: func(ctx context.Context, query, report interface{}) error {
				return storage.ErrCertificateNotFound
			},
		}

		sut := handler.NewCertificateIssuer(commandHandler, queryRunner)

		// act
		_, err := sut.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, storage.ErrCertificateNotFound)
	})

	t.Run("an error occurred when getting a certificate, internal system error returned", func(t *testing.T) {
		t.Parallel()
		// arrange
		certificateID := gofakeit.UUID()

		rq := &req.IssueCertificate{
			CertificateID: certificateID,
			ArtistID:      gofakeit.UUID(),
			Title:         gofakeit.Sentence(5),
			ArtworkType:   "painting",
		}

		commandHandler := CommandHandlerSpy{
			Fn: func(ctx context.Context, command interface{}) error {
				return nil
			},
		}

		queryRunner := QueryRunnerSpy{
			Fn: func(ctx context.Context, query, report interface{}) error {
				return errSomeBadThingHappened
			},
		}

		sut := handler.NewCertificateIssuer(commandHandler, queryRunner)

		// act
		_, err := sut.Handle(context.Background(), rq)

		// assert
		assertCause(t, err, errSomeBadThingHappened)
	})

	t.Run("certificate issued successfully, valid response returned", func(t *testing.T) {
		t.Parallel()

		// arrange
		var want *resp.IssueCertificate
		gofakeit.Struct(&want)

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

		commandHandler := CommandHandlerSpy{
			Fn: func(ctx context.Context, command interface{}) error {
				assert.Equals(t, c, command)

				return nil
			},
		}

		queryRunner := QueryRunnerSpy{
			Fn: func(ctx context.Context, query, rep interface{}) error {
				r, ok := rep.(*report.Certificate)
				assert.True(t, ok)

				*r = report.Certificate{
					ArtistName:  want.ArtistName,
					ArtworkType: want.ArtworkType,
					Title:       want.Title,
				}

				return nil
			},
		}

		sut := handler.NewCertificateIssuer(commandHandler, queryRunner)

		// act
		res, err := sut.Handle(context.Background(), rq)

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
