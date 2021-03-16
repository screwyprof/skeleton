package resp

type ViewCertificate struct {
	CertificateID string `json:"certificate_id"`
	Title         string `json:"title"`
	ArtworkType   string `json:"artwork_type"`
	ArtistName    string `json:"artist_name"`
}
