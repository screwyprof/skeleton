//go:build integration
// +build integration

package issuecert_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/screwyprof/golibs/cmdhandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/skeleton/pkg/cert/command"
	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
)

func TestIssueCertificate(t *testing.T) {
	// arrange
	certStorage := &CertStorageSpy{}

	certIssuerHandler, err := initCommandHandler(certStorage)
	require.NoError(t, err)

	// act
	c := command.IssueCertificate{
		ID:          gofakeit.UUID(),
		Title:       gofakeit.Sentence(5),
		ArtistID:    gofakeit.UUID(),
		ArtworkType: "painting",
	}
	err = certIssuerHandler.Handle(context.Background(), c)

	// assert
	require.NoError(t, err)
	assert.True(t, certStorage.wasCalled)
}

func initCommandHandler(certStorage issuecert.CertStorage) (*cmdhandler.Dispatcher, error) {
	certIssuer := issuecert.NewHandler(certStorage)
	certIssuerHandler, err := cmdhandler.Adapt(certIssuer.Handle)
	if err != nil {
		return nil, err
	}
	d := cmdhandler.NewDispatcher()
	d.RegisterCommandHandler("IssueCertificate", certIssuerHandler)

	return d, nil
}
