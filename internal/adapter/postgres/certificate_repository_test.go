package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v4"
	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/go-rel/reltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/skeleton/cert/report"
	"github.com/screwyprof/skeleton/cert/usecase/issuecert"
	"github.com/screwyprof/skeleton/cert/usecase/storage"
	"github.com/screwyprof/skeleton/internal/adapter/postgres"
	"github.com/screwyprof/skeleton/internal/adapter/postgres/model"
)

var ErrUnknownFailure = errors.New("unknown failure")

func TestCertificateCertificateByID(t *testing.T) {
	t.Parallel()

	t.Run("it returns certificate by id", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := report.Certificate{
			ID:          gofakeit.UUID(),
			Title:       gofakeit.Sentence(3),
			ArtistName:  "Shepard Fairey",
			ArtworkType: gofakeit.Word(),
		}

		row := model.Certificate{
			CertificateID: want.ID,
			Title:         want.Title,
			ArtistID:      want.ArtistName,
			ArtworkType:   want.ArtworkType,
		}

		repo := reltest.New()
		repo.ExpectFind(where.Eq("certificate_id", want.ID)).Result(row)

		sut := postgres.NewCertificateRepository(repo)

		// act
		c, err := sut.CertificateByID(context.Background(), want.ID)

		// assert
		require.NoError(t, err)
		assert.Equal(t, want, c)
		repo.AssertExpectations(t)
	})

	t.Run("it returns an error if certificate is not found", func(t *testing.T) {
		t.Parallel()

		// arrange
		certificateID := gofakeit.UUID()

		repo := reltest.New()
		repo.ExpectFind(where.Eq("certificate_id", certificateID)).Error(rel.ErrNotFound)

		sut := postgres.NewCertificateRepository(repo)

		// act
		_, err := sut.CertificateByID(context.Background(), certificateID)

		// assert
		assert.ErrorIs(t, err, storage.ErrCertificateNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("it returns an error if it cannot fetch a certificate", func(t *testing.T) {
		t.Parallel()

		// arrange
		certificateID := gofakeit.UUID()

		repo := reltest.New()
		repo.ExpectFind(where.Eq("certificate_id", certificateID)).Error(ErrUnknownFailure)

		sut := postgres.NewCertificateRepository(repo)

		// act
		_, err := sut.CertificateByID(context.Background(), certificateID)

		// assert
		assert.ErrorIs(t, err, ErrUnknownFailure)
		repo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	t.Parallel()

	t.Run("it stores a certificate", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := &issuecert.Certificate{
			ID:          gofakeit.UUID(),
			Title:       gofakeit.Sentence(3),
			ArtistID:    "Shepard Fairey",
			ArtworkType: gofakeit.Word(),
		}

		row := model.Certificate{
			CertificateID: want.ID,
			Title:         want.Title,
			ArtistID:      want.ArtistID,
			ArtworkType:   want.ArtworkType,
		}

		repo := reltest.New()
		repo.ExpectTransaction(func(repo *reltest.Repository) {
			repo.ExpectInsert().For(&row)
		})

		sut := postgres.NewCertificateRepository(repo)

		// act
		err := sut.Store(context.Background(), want)

		// assert
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("it returns an error if certificate already exists", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := &issuecert.Certificate{
			ID:          gofakeit.UUID(),
			Title:       gofakeit.Sentence(3),
			ArtistID:    "Shepard Fairey",
			ArtworkType: gofakeit.Word(),
		}

		row := model.Certificate{
			CertificateID: want.ID,
			Title:         want.Title,
			ArtistID:      want.ArtistID,
			ArtworkType:   want.ArtworkType,
		}

		repo := reltest.New()
		repo.ExpectTransaction(func(repo *reltest.Repository) {
			repo.ExpectInsert().For(&row).Error(rel.ErrUniqueConstraint)
		})

		sut := postgres.NewCertificateRepository(repo)

		// act
		err := sut.Store(context.Background(), want)

		// assert
		assert.ErrorIs(t, err, storage.ErrCertificateAlreadyExists)
		repo.AssertExpectations(t)
	})

	t.Run("it returns an error if it cannot store a certificate", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := &issuecert.Certificate{
			ID:          gofakeit.UUID(),
			Title:       gofakeit.Sentence(3),
			ArtistID:    "Shepard Fairey",
			ArtworkType: gofakeit.Word(),
		}

		row := model.Certificate{
			CertificateID: want.ID,
			Title:         want.Title,
			ArtistID:      want.ArtistID,
			ArtworkType:   want.ArtworkType,
		}

		repo := reltest.New()
		repo.ExpectTransaction(func(repo *reltest.Repository) {
			repo.ExpectInsert().For(&row).Error(ErrUnknownFailure)
		})

		sut := postgres.NewCertificateRepository(repo)

		// act
		err := sut.Store(context.Background(), want)

		// assert
		assert.ErrorIs(t, err, ErrUnknownFailure)
		repo.AssertExpectations(t)
	})
}
