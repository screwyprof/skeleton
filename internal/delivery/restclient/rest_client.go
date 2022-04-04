package restclient

import (
	"encoding/json"
	stdErrors "errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/screwyprof/skeleton/internal/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/delivery/rest/resp"
)

var ErrUnexpectedHTTPStatus = stdErrors.New("unexpected http status")

type RESTClient struct {
	httpClient *resty.Client
}

func New(baseURL string) *RESTClient {
	r := resty.New()
	r.SetBaseURL(baseURL)

	return &RESTClient{httpClient: r}
}

func (c *RESTClient) IssueCertificate(certificate req.IssueCertificate) (*resp.ViewCertificate, error) {
	jsonBody, err := json.Marshal(certificate)
	if err != nil {
		return nil, errors.Wrap(err, "cannot issue certificate")
	}

	r, err := c.httpClient.R().
		SetBody(jsonBody).
		SetResult(&resp.ViewCertificate{}).
		Post("/api/v1/certificates")
	if err != nil {
		return nil, errors.Wrap(err, "cannot issue certificate")
	}

	if r.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("%w: unexpected http status %d", ErrUnexpectedHTTPStatus, r.StatusCode())
	}

	res := r.Result().(*resp.ViewCertificate) //nolint:forcetypeassert

	return res, nil
}

func (c *RESTClient) ViewCertificate(certificateID string) (*resp.ViewCertificate, error) {
	r, err := c.httpClient.R().
		SetPathParams(map[string]string{
			"certificate_id": certificateID,
		}).SetResult(&resp.ViewCertificate{}).
		Get("/api/v1/certificates/{certificate_id}")
	if err != nil {
		return nil, errors.Wrap(err, "cannot retrieve certificate")
	}

	if r.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%w: unexpected http status %d", ErrUnexpectedHTTPStatus, r.StatusCode())
	}

	res := r.Result().(*resp.ViewCertificate) //nolint:forcetypeassert

	return res, nil
}
