package req

import (
	"github.com/gin-gonic/gin"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/apierr"
)

type ViewCertificate struct {
	CertificateID string `uri:"certificate_id" binding:"required,uuid"`
}

func (r *ViewCertificate) Bind(c *gin.Context) error {
	if err := c.ShouldBindUri(r); err != nil {
		return apierr.Wrap(err, apierr.BadRequest, nil)
	}

	return nil
}
