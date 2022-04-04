package viewcert_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/screwyprof/golibs/assert"
	. "github.com/screwyprof/golibs/queryer/testdsl"

	"github.com/screwyprof/skeleton/cert/query"
	"github.com/screwyprof/skeleton/cert/report"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
	"github.com/screwyprof/skeleton/cert/usecase/viewcert"
)

func TestViewCertificate(t *testing.T) {
	t.Parallel()

	t.Run("non existent certificate id given, not found error returned", func(t *testing.T) {
		t.Parallel()

		certificateID := gofakeit.UUID()

		reporter := CertReporterSpy{
			Fn: func(ctx context.Context, id string) (report.Certificate, error) {
				return report.Certificate{}, storage.ErrCertificateNotFound
			},
		}

		sut := viewcert.NewQueryer(reporter).ViewCertificate

		Test(t)(
			Given("ViewCertificate", sut),
			When(context.Background(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}),
			ThenFailWith(storage.ErrCertificateNotFound),
		)
	})

	t.Run("an existing certificate id given, certificate info returned", func(t *testing.T) {
		t.Parallel()

		var want *report.Certificate
		gofakeit.Struct(&want)

		certificateID := gofakeit.UUID()
		want.ID = certificateID

		reporter := CertReporterSpy{
			Fn: func(ctx context.Context, id string) (report.Certificate, error) {
				assert.Equals(t, certificateID, id)

				return *want, nil
			},
		}

		sut := viewcert.NewQueryer(reporter).ViewCertificate

		Test(t)(
			Given("ViewCertificate", sut),
			When(context.Background(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}),
			Then(want),
		)
	})
}

type CertReporterSpy struct {
	Fn func(ctx context.Context, certificateID string) (report.Certificate, error)
}

func (s CertReporterSpy) CertificateByID(ctx context.Context, certificateID string) (report.Certificate, error) {
	return s.Fn(ctx, certificateID)
}
