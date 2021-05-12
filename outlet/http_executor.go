package outlet

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/funcas/cgs/connector"

	"github.com/funcas/cgs/charset"

	"github.com/funcas/cgs/message"
)

var _ Executor = (*HttpExecutor)(nil)

type HttpExecutor struct {
	BaseExecutor
}

func NewHttpExecutor(conn connector.Connector) *HttpExecutor {
	return &HttpExecutor{
		BaseExecutor{
			connector: conn,
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
	req, _ := http.NewRequest(conn.Method(), conn.Url(), nil)
	req.Header.Set("content-type", conn.ContentType())
	if conn.UserAgent() != "" {
		req.Header.Set("user-agent", conn.UserAgent())
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	dst, _ := charset.ToUTF8(charset.Charset(conn.Lang()), string(body))
	msg.OriData = dst

}
