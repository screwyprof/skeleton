//go:build integration
// +build integration

package viewcert_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/screwyprof/golibs/queryer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/skeleton/cert/query"
	"github.com/screwyprof/skeleton/cert/report"
	"github.com/screwyprof/skeleton/cert/usecase/viewcert"
)

func TestViewCertificate(t *testing.T) {
	// arrange
	certificateID := gofakeit.UUID()
	want := report.Certificate{
		Title:       gofakeit.Sentence(5),
		ArtistName:  gofakeit.Name(),
		ArtworkType: "painting",
	}

	certificateViewer, err := initQueryer(CertReporterWithValidCertificateStub{Want: want})
	require.NoError(t, err)

	// act
	var r report.Certificate
	err = certificateViewer.RunQuery(context.Background(), query.ViewCertificate{ID: certificateID}, &r)

	// assert
	require.NoError(t, err)
	assert.Equal(t, want, r)
}

func initQueryer(reporter viewcert.CertReporter) (*queryer.Dispatcher, error) {
	certViewer := viewcert.NewQueryer(reporter)
	certViewerUseCaseRunner, err := queryer.Adapt(certViewer.ViewCertificate)
	if err != nil {
		return nil, err
	}
	d := queryer.NewDispatcher()
	d.RegisterQueryRunner("ViewCertificate", certViewerUseCaseRunner)

	return d, nil
}
