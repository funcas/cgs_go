package connector

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/valyala/fastjson"
)

var _ Connector = (*WebserviceConnector)(nil)

type WebserviceConnector struct {
	BaseConnector
	wsAddress       string
	namespace       string
	methodName      string
	soapAction      string
	methodNamespace string
	paramNames      string
	timeout         int
	proxy           string
}

func NewWebserviceConnector(vo *fastjson.Value) *WebserviceConnector {
	return &WebserviceConnector{
		wsAddress:       string(vo.GetStringBytes("wsAddress")),
		namespace:       string(vo.GetStringBytes("namespace")),
		methodName:      string(vo.GetStringBytes("methodName")),
		soapAction:      string(vo.GetStringBytes("soapAction")),
		methodNamespace: string(vo.GetStringBytes("methodNamespace")),
		paramNames:      string(vo.GetStringBytes("paramNames")),
		timeout:         vo.GetInt("timeout"),
		proxy:           string(vo.GetStringBytes("proxy")),
		BaseConnector: BaseConnector{
			name:    string(vo.GetStringBytes("name")),
			enabled: vo.GetBool("enabled"),
		},
	}
}

func (w WebserviceConnector) WsAddress() string {
	return w.wsAddress
}

func (w WebserviceConnector) Namespace() string {
	return w.namespace
}

func (w WebserviceConnector) MethodName() string {
	return w.methodName
}

func (w WebserviceConnector) SoapAction() string {
	return w.soapAction
}

func (w WebserviceConnector) MethodNamespace() string {
	return w.methodNamespace
}

func (w WebserviceConnector) ParamNames() string {
	return w.paramNames
}

func (w WebserviceConnector) GetHttpClient() (client *http.Client, err error) {
	if !w.enabled {
		err = errors.New("connector has been disabled")
		return
	}

	client = &http.Client{
		Timeout: time.Duration(w.timeout) * time.Second,
	}
	// 若设置了代码地址，则使用代理
	if w.proxy != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(w.proxy)
		}
		client.Transport = &http.Transport{
			Proxy: proxy,
		}
	}

	return
}
