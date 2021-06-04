package connector

import (
	"errors"
	"net/http"
	"time"

	"github.com/valyala/fastjson"
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
	lang        string
	header      string
}

func NewHttpConnector(vo *fastjson.Value) *HttpConnector {
	return &HttpConnector{
		url:         string(vo.GetStringBytes("url")),
		method:      string(vo.GetStringBytes("method")),
		queryParams: string(vo.GetStringBytes("queryParams")),
		contentType: string(vo.GetStringBytes("contentType")),
		userAgent:   string(vo.GetStringBytes("userAgent")),
		timeout:     vo.GetInt("timeout"),
		lang:        string(vo.GetStringBytes("lang")),
		header:      string(vo.GetStringBytes("header")),
		BaseConnector: BaseConnector{
			name:    string(vo.GetStringBytes("name")),
			enabled: vo.GetBool("enabled"),
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

func (h HttpConnector) Lang() string {
	return h.lang
}

func (h HttpConnector) Header() string {
	return h.header
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
