package storage

import "errors"

var (
	ErrCertificateNotFound      = errors.New("cannot store certificate")
	ErrCannotStoreCertificate   = errors.New("cannot store certificate")
	ErrCertificateAlreadyExists = errors.New("certificate already exists")
)
