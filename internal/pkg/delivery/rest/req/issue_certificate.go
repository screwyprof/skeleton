package req

import (
	"github.com/gin-gonic/gin"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/apierr"
)

type IssueCertificate struct {
	CertificateID string `binding:"required,uuid" json:"certificate_id"`
	ArtistID      string `binding:"required,uuid" json:"artist_id"`
	Title         string `binding:"required" json:"title"`
	ArtworkType   string `binding:"required" json:"artwork_type"`
}

func (r *IssueCertificate) Bind(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return apierr.Wrap(err, apierr.BadRequest, nil)
	}

	return nil
}
