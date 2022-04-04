package modqueryer

import (
	"github.com/screwyprof/golibs/queryer"
	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/cert/usecase/viewcert"
)

var Module = fx.Provide(
	New,
)

func New(reporter viewcert.CertReporter) queryer.QueryRunner {
	certViewer := viewcert.NewQueryer(reporter)
	certViewerUseCaseRunner := queryer.MustAdapt(certViewer.ViewCertificate)

	d := queryer.NewDispatcher()
	d.RegisterQueryRunner("ViewCertificate", certViewerUseCaseRunner)

	return d
}
