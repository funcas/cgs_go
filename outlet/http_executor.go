package outlet

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/funcas/cgs/analysis"

	"github.com/funcas/cgs/tpl"

	"github.com/funcas/cgs/connector"

	"github.com/funcas/cgs/charset"

	"github.com/funcas/cgs/message"
)

var _ Executor = (*HttpExecutor)(nil)

type HttpExecutor struct {
	BaseExecutor
}

func NewHttpExecutor(conn connector.Connector, tpl tpl.TemplateService, analysis analysis.Analyser) *HttpExecutor {
	return &HttpExecutor{
		BaseExecutor{
			connector: conn,
			template:  tpl,
			analysis:  analysis,
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
		log.Println(err.Error())
	}
	var req *http.Request
	url := buildUrlParams(conn, msg)
	if conn.Method() == "POST" {
		reqBody := exe.BaseExecutor.template.GetTemplateFromMessage(*msg)
		log.Printf("template content is >>> %s <<<\n", reqBody)
		req, _ = http.NewRequest(conn.Method(), url, bytes.NewBufferString(reqBody))
	} else {
		req, _ = http.NewRequest(conn.Method(), url, nil)
	}

	req.Header.Set("content-type", conn.ContentType())
	if conn.UserAgent() != "" {
		req.Header.Set("user-agent", conn.UserAgent())
	}

	// add headers to Http Header
	header := conn.Header()
	if header != "" {
		headerKeys := strings.Split(header, ",")
		for _, h := range headerKeys {
			req.Header.Set(h, msg.Params[h])
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		msg.OriData = err.Error()
		log.Println(err.Error())
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
	if exe.analysis != nil {
		msg.Data = exe.analysis.AnalysisResult(dst, msg.TransCode)
	}

}

// build url parameters into [key=val&key2=val2] from connector configuration
func buildUrlParams(conn *connector.HttpConnector, msg *message.Message) string {
	params := strings.Split(conn.QueryParams(), ",")
	var bf bytes.Buffer
	for _, param := range params {
		bf.WriteString("&")
		bf.WriteString(param)
		bf.WriteString("=")
		bf.WriteString(msg.Params[param])
	}
	url := conn.Url()
	if strings.ContainsAny(url, "?") {
		return url + bf.String()
	} else {
		return url + "?" + bf.String()[1:]
	}
}
