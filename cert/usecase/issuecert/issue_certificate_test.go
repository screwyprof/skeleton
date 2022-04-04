package issuecert_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/golang/mock/gomock"
	. "github.com/screwyprof/golibs/cmdhandler/testdsl"

	"github.com/screwyprof/skeleton/cert/command"
	"github.com/screwyprof/skeleton/cert/mock"
	"github.com/screwyprof/skeleton/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
)

func TestIssueCertificate(t *testing.T) {
	t.Parallel()

	t.Run("valid command given, certificate created", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()
		artistID := gofakeit.UUID()
		title := gofakeit.Sentence(5)

		want := &issuecert.Certificate{
			ID:          certificateID,
			ArtistID:    artistID,
			ArtworkType: "painting",
			Title:       title,
		}

		certStorage := createCertStorage(ctrl, want)
		handler := issuecert.NewHandler(certStorage).Handle

		Test(t)(
			Given("IssueCertificate", handler),
			When(context.Background(), command.IssueCertificate{
				ID:          certificateID,
				ArtistID:    artistID,
				ArtworkType: "painting",
				Title:       title,
			}),
			ThenOk(),
		)
	})

	t.Run("it returns an error if it fails to store a certificate", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		certificateID := gofakeit.UUID()
		artistID := gofakeit.UUID()
		title := gofakeit.Sentence(5)

		certStorage := createFaultyCertStorage(ctrl, storage.ErrCannotStoreCertificate)
		handler := issuecert.NewHandler(certStorage).Handle

		Test(t)(
			Given("IssueCertificate", handler),
			When(context.Background(), command.IssueCertificate{
				ID:          certificateID,
				ArtistID:    artistID,
				ArtworkType: "painting",
				Title:       title,
			}),
			ThenFailWith(storage.ErrCannotStoreCertificate),
		)
	})
}

func createCertStorage(ctrl *gomock.Controller, want *issuecert.Certificate) *mock.MockCertStorage {
	certStorage := mock.NewMockCertStorage(ctrl)
	certStorage.EXPECT().
		Store(context.Background(), want).
		Return(nil)

	return certStorage
}

func createFaultyCertStorage(ctrl *gomock.Controller, err error) *mock.MockCertStorage {
	certStorage := mock.NewMockCertStorage(ctrl)
	certStorage.EXPECT().
		Store(context.Background(), gomock.Any()).
		Return(err)

	return certStorage
}
