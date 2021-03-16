package restclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/req"
	"github.com/screwyprof/skeleton/internal/pkg/delivery/rest/resp"
)

type RESTClient struct {
	httpClient *resty.Client
}

func New(baseURL string) *RESTClient {
	r := resty.New()
	r.SetHostURL(baseURL)

	return &RESTClient{httpClient: r}
}

func (c *RESTClient) IssueCertificate(certificate req.IssueCertificate) (*resp.ViewCertificate, error) {
	jsonBody, err := c.jsonStringFor(certificate)
	if err != nil {
		return nil, err
	}

	r, err := c.httpClient.R().
		SetBody(jsonBody).
		SetResult(&resp.ViewCertificate{}).
		Post("/api/v1/certificates")
	if err != nil {
		return nil, err
	}

	if r.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("unexpected http status %d", r.StatusCode())
	}

	res := r.Result().(*resp.ViewCertificate)

	return res, nil
}

func (c *RESTClient) ViewCertificate(certificateID string) (*resp.ViewCertificate, error) {
	r, err := c.httpClient.R().
		SetPathParams(map[string]string{
			"certificate_id": certificateID,
		}).SetResult(&resp.ViewCertificate{}).
		Get("/api/v1/certificates/{certificate_id}")
	if err != nil {
		return nil, err
	}

	if r.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status %d", r.StatusCode())
	}

	res := r.Result().(*resp.ViewCertificate)

	return res, nil
}

func (c *RESTClient) jsonStringFor(object interface{}) (string, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
