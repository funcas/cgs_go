package connector

import (
	"errors"
	"net/http"
	"time"

	"github.com/funcas/cgs/model"
)

var _ Connector = (*HttpConnector)(nil)

type HttpConnector struct {
	BaseConnector
	url         string
	method      string
	queryParams string
	contentType string
	userAgent   string
	timeout     int
}

func NewHttpConnector(vo model.HttpConnectorVO) *HttpConnector {
	return &HttpConnector{
		url:         vo.URL,
		method:      vo.Method,
		queryParams: vo.QueryParams,
		contentType: vo.ContentType,
		userAgent:   vo.UserAgent,
		timeout:     vo.Timeout,
		BaseConnector: BaseConnector{
			name:    vo.Name,
			enabled: vo.Enabled,
		},
	}
}

func (h *HttpConnector) SetUrl(url string) {
	h.url = url
}

func (h HttpConnector) Url() string {
	return h.url
}

func (h HttpConnector) Method() string {
	return h.method
}

func (h HttpConnector) QueryParams() string {
	return h.queryParams
}

func (h HttpConnector) UserAgent() string {
	return h.userAgent
}

func (h HttpConnector) Timeout() int {
	return h.timeout
}

func (h HttpConnector) ContentType() string {
	return h.contentType
}

func (h HttpConnector) GetHttpClient() (client *http.Client, err error) {
	if !h.enabled {
		err = errors.New("connector has been disabled")
		return
	}
	client = &http.Client{
		Timeout: time.Duration(h.timeout) * time.Second,
	}

	return
}
