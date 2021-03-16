package resp

import "net/http"

type IssueCertificate struct {
	CertificateID string `json:"certificate_id"`
	Title         string `json:"title"`
	ArtworkType   string `json:"artwork_type"`
	ArtistName    string `json:"artist_name"`
}

func (r *IssueCertificate) Status() int {
	return http.StatusCreated
}
