package outlet

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/funcas/cgs/tpl"

	"github.com/funcas/cgs/connector"

	"github.com/funcas/cgs/charset"

	"github.com/funcas/cgs/message"
)

var _ Executor = (*HttpExecutor)(nil)

type HttpExecutor struct {
	BaseExecutor
}

func NewHttpExecutor(conn connector.Connector, tpl tpl.TemplateService) *HttpExecutor {
	return &HttpExecutor{
		BaseExecutor{
			connector: conn,
			template:  tpl,
		},
	}
}

func (exe HttpExecutor) CanExecute(transCode string) bool {
	return true
}

func (exe HttpExecutor) Execute(msg *message.Message) {

	conn := exe.BaseExecutor.connector.(*connector.HttpConnector)
	client, err := conn.GetHttpClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	var req *http.Request
	if conn.Method() == "POST" {
		reqBody := exe.BaseExecutor.template.GetTemplateFromMessage(*msg)
		log.Printf("template content is >>> %s <<<", reqBody)
		req, _ = http.NewRequest(conn.Method(), conn.Url(), bytes.NewBufferString(reqBody))
	} else {
		req, _ = http.NewRequest(conn.Method(), conn.Url(), nil)
	}

	req.Header.Set("content-type", conn.ContentType())
	if conn.UserAgent() != "" {
		req.Header.Set("user-agent", conn.UserAgent())
	}
	resp, err := client.Do(req)
	if err != nil {
		msg.OriData = err.Error()
		log.Fatal(err.Error())
		return
	}
	if resp.StatusCode != 200 {
		msg.OriData = "request error"
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	dst, _ := charset.ToUTF8(charset.Charset(conn.Lang()), string(body))
	log.Printf("return msg >>> %s <<<", dst)
	msg.OriData = dst

}

// build url parameters into [key=val&key2=val2] from connector configuration
func buildUrlParams(msg *message.Message) string {
	return ""
}
