package model

type Certificate struct {
	CertificateID string `pg:"certificate_id,pk"`
	ArtistID      string
	Title         string
	ArtworkType   string
}
