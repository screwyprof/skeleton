package viewcert_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/golang/mock/gomock"
	. "github.com/screwyprof/golibs/queryer/testdsl"

	"github.com/screwyprof/skeleton/pkg/cert/mock"
	"github.com/screwyprof/skeleton/pkg/cert/query"
	"github.com/screwyprof/skeleton/pkg/cert/report"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/storage"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/viewcert"
)

func TestViewCertificate(t *testing.T) {
	t.Parallel()

	t.Run("non existent certificate id given, not found error returned", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()

		reporter := createErrCertReporter(ctrl, certificateID, storage.ErrCertificateNotFound)
		concreteQueryRunner := viewcert.NewQueryer(reporter).ViewCertificate

		Test(t)(
			Given("ViewCertificate", concreteQueryRunner),
			When(context.Background(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}),
			ThenFailWith(storage.ErrCertificateNotFound),
		)
	})

	t.Run("an existing certificate id given, certificate info returned", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()
		want := report.Certificate{
			ArtistName:  gofakeit.Name(),
			ArtworkType: gofakeit.BuzzWord(),
			Title:       gofakeit.Sentence(5),
		}

		reporter := createCertReporter(ctrl, certificateID, want)
		concreteQueryRunner := viewcert.NewQueryer(reporter).ViewCertificate

		Test(t)(
			Given("ViewCertificate", concreteQueryRunner),
			When(context.Background(), query.ViewCertificate{ID: certificateID}, &report.Certificate{}),
			Then(&want),
		)
	})
}

func createCertReporter(ctrl *gomock.Controller, certificateID string, want report.Certificate) *mock.MockCertReporter {
	certReporter := mock.NewMockCertReporter(ctrl)
	certReporter.EXPECT().
		CertificateByID(context.Background(), certificateID).
		Return(want, nil)

	return certReporter
}

func createErrCertReporter(ctrl *gomock.Controller, certificateID string, want error) *mock.MockCertReporter {
	certReporter := mock.NewMockCertReporter(ctrl)
	certReporter.EXPECT().
		CertificateByID(context.Background(), certificateID).
		Return(report.Certificate{}, want)

	return certReporter
}
