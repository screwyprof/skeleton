package modstorerep

import (
	"github.com/go-rel/rel"
	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/internal/adapter/postgres"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/viewcert"
)

var Module = fx.Provide(
	NewStorage,
	NewReporter,
)

func NewStorage(repo rel.Repository) issuecert.CertStorage {
	return postgres.NewCertificateRepository(repo)
}

func NewReporter(repo rel.Repository) viewcert.CertReporter {
	return postgres.NewCertificateRepository(repo)
}
