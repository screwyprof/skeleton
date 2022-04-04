package issuecert_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	. "github.com/screwyprof/golibs/cmdhandler/testdsl"

	"github.com/screwyprof/skeleton/cert/command"
	"github.com/screwyprof/skeleton/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
)

func TestIssueCertificate(t *testing.T) {
	t.Parallel()

	t.Run("valid command given, certificate created", func(t *testing.T) {
		t.Parallel()

		var c command.IssueCertificate
		gofakeit.Struct(&c)

		want := issuecert.Certificate{
			ID:          c.ID,
			Title:       c.Title,
			ArtistID:    c.ArtistID,
			ArtworkType: c.ArtworkType,
		}

		certStorage := CertStorageSpy{
			Fn: func(ctx context.Context, c *issuecert.Certificate) error {
				*c = want

				return nil
			},
		}

		sut := issuecert.NewHandler(certStorage).Handle

		Test(t)(
			Given("IssueCertificate", sut),
			When(context.Background(), c),
			ThenOk(),
		)
	})

	t.Run("it returns an error if it fails to store a certificate", func(t *testing.T) {
		t.Parallel()

		var c command.IssueCertificate
		gofakeit.Struct(&c)

		certStorage := CertStorageSpy{
			Fn: func(ctx context.Context, c *issuecert.Certificate) error {
				return storage.ErrCannotStoreCertificate
			},
		}

		sut := issuecert.NewHandler(certStorage).Handle

		Test(t)(
			Given("IssueCertificate", sut),
			When(context.Background(), c),
			ThenFailWith(storage.ErrCannotStoreCertificate),
		)
	})
}

type CertStorageSpy struct {
	Fn func(ctx context.Context, c *issuecert.Certificate) error
}

func (s CertStorageSpy) Store(ctx context.Context, c *issuecert.Certificate) error {
	return s.Fn(ctx, c)
}
