package issuecert

import (
	"context"

	"github.com/screwyprof/skeleton/pkg/cert/command"
)

type CertIssuer struct {
	storage CertStorage
}

type CertStorage interface {
	Store(context.Context, *Certificate) error
}

func NewHandler(storage CertStorage) *CertIssuer {
	return &CertIssuer{storage: storage}
}

func (i CertIssuer) Handle(ctx context.Context, q command.IssueCertificate) error {
	cert := &Certificate{
		ID:          q.ID,
		ArtistID:    q.ArtistID,
		ArtworkType: q.ArtworkType,
		Title:       q.Title,
	}

	return i.storage.Store(ctx, cert)
}
