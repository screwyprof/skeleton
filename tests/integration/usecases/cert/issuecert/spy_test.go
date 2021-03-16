//+build integration

package issuecert_test

import (
	"context"

	"github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
)

type CertStorageSpy struct {
	wasCalled bool
}

func (s *CertStorageSpy) Store(_ context.Context, _ *issuecert.Certificate) error {
	s.wasCalled = true
	return nil
}
