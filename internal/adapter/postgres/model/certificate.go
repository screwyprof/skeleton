package model

type Certificate struct {
	CertificateID string `db:"certificate_id,uuid,primary"`
	ArtistID      string
	Title         string
	ArtworkType   string
}
