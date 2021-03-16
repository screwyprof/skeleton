//+build integration

package viewcert_test

import (
	"context"

	"github.com/screwyprof/skeleton/pkg/cert/report"
)

type CertReporterWithValidCertificateStub struct {
	Want report.Certificate
}

func (s CertReporterWithValidCertificateStub) CertificateByID(_ context.Context, _ string) (report.Certificate, error) {
	return s.Want, nil
}
