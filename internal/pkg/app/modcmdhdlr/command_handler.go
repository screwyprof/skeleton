package modcmdhdlr

import (
	"github.com/screwyprof/golibs/cmdhandler"
	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
)

var Module = fx.Provide(
	New,
)

func New(certStorage issuecert.CertStorage) cmdhandler.CommandHandler {
	certIssuer := issuecert.NewHandler(certStorage)
	certIssuerHandler := cmdhandler.MustAdapt(certIssuer.Handle)

	d := cmdhandler.NewDispatcher()
	d.RegisterCommandHandler("IssueCertificate", certIssuerHandler)

	return d
}
