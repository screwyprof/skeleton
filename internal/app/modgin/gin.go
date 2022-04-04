package modgin

import (
	"net/http"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/screwyprof/golibs/adaptor"
	"github.com/screwyprof/golibs/cmdhandler"
	"github.com/screwyprof/golibs/gin/middleware/ctxtags"
	"github.com/screwyprof/golibs/gin/middleware/ctxzap"
	"github.com/screwyprof/golibs/gin/middleware/errors"
	"github.com/screwyprof/golibs/gin/renderer"
	"github.com/screwyprof/golibs/queryer"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/screwyprof/skeleton/internal/delivery/rest/handler"
	"github.com/screwyprof/skeleton/internal/delivery/rest/req"
)

var Module = fx.Provide(
	NewGin,
	NewHTTPHandler,
)

func NewGin(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	e := gin.New()
	e.Use(ginzap.RecoveryWithZap(logger, true))
	e.Use(errors.ErrorHandler())
	e.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	// e.Use(ctxtx.CtxTX(db))
	e.Use(ctxtags.CtxTags(ctxtags.WithFieldExtractor(ctxtags.RequestID)))
	e.Use(ctxzap.CtxZap(logger, time.RFC3339, true))

	return e
}

func NewHTTPHandler(
	mux *gin.Engine,
	commandHandler cmdhandler.CommandHandler,
	queryRunner queryer.QueryRunner,
) http.Handler {
	certViewer := adaptor.MustAdapt(handler.NewCertificateViewer(queryRunner).Handle)
	certIssuer := adaptor.MustAdapt(handler.NewCertificateIssuer(commandHandler, queryRunner).Handle)

	r := renderer.NewGinRenderer()
	r.Register(&req.ViewCertificate{}, certViewer)
	r.Register(&req.IssueCertificate{}, certIssuer)

	v1 := mux.Group("/api/v1")
	{
		certificates := v1.Group("/certificates")
		{
			certificates.GET("/:certificate_id", r.MustAdapt(certViewer))
			certificates.POST("", r.MustAdapt(certIssuer))
		}
	}

	return mux
}
