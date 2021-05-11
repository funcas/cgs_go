package executor

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"

	"golang.org/x/text/transform"

	"github.com/funcas/cgs/connector"
	"github.com/funcas/cgs/message"
)

var _ Executor = (*HttpExecutor)(nil)

type HttpExecutor struct {
	httpConnector *connector.HttpConnector
}

func (exe *HttpExecutor) SetHttpConnector(conn *connector.HttpConnector) {
	exe.httpConnector = conn
}

func (exe HttpExecutor) CanExecute(transCode string) bool {
	return true
}

func (exe HttpExecutor) Execute(msg message.Message) {
	conn := exe.httpConnector
	client, err := conn.GetHttpClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	req, _ := http.NewRequest(conn.Method(), conn.Url(), nil)
	req.Header.Set("content-type", conn.ContentType())
	req.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	if conn.UserAgent() != "" {
		req.Header.Set("user-agent", conn.UserAgent())
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := transform.NewReader(bytes.NewReader(body), simplifiedchinese.GBK.NewEncoder())
	ret, _ := ioutil.ReadAll(data)
	log.Println(string(ret))

}
