package modstorerep

import (
	"github.com/go-pg/pg/v9"
	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/internal/pkg/adapter/postgres"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/viewcert"
)

var Module = fx.Provide(
	NewStorage,
	NewReporter,
)

func NewStorage(db *pg.DB) issuecert.CertStorage {
	return postgres.NewCertificateRepository(postgres.NewCtxTxRunner(db))
}

func NewReporter(db *pg.DB) viewcert.CertReporter {
	return postgres.NewCertificateRepository(postgres.NewCtxTxRunner(db))
}
