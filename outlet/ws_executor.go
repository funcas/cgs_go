// webservice executor
package outlet

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/funcas/cgs/analysis"
	"github.com/funcas/cgs/tpl"

	"github.com/funcas/cgs/connector"
	"github.com/funcas/cgs/message"
)

var _ Executor = (*WSExecutor)(nil)

type WSExecutor struct {
	BaseExecutor
}

func NewWSExecutor(conn connector.Connector, tpl tpl.TemplateService, analysis analysis.Analyser) *WSExecutor {
	return &WSExecutor{
		BaseExecutor{
			connector: conn,
			template:  tpl,
			analysis:  analysis,
		},
	}
}

func (exe WSExecutor) CanExecute(transCode string) bool {
	return true
}

func (exe WSExecutor) Execute(msg *message.Message) {
	conn := exe.BaseExecutor.connector.(*connector.WebserviceConnector)
	client, err := conn.GetHttpClient()
	if err != nil {
		log.Println(err.Error())
	}
	var req *http.Request
	bodyXml := prepareSOAPMsg(msg, exe, conn)
	log.Println("sending xml: ")
	log.Println(bodyXml)
	req, _ = http.NewRequest("POST", conn.WsAddress(), bytes.NewBufferString(bodyXml))
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Content-Length", strconv.Itoa(len(bodyXml)))
	if conn.SoapAction() != "" {
		req.Header.Set("SOAPAction", conn.SoapAction())
	} else {
		req.Header.Set("SOAPAction", conn.MethodName())
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
	}
	ret, _ := ioutil.ReadAll(resp.Body)
	msg.OriData = string(ret)
	if resp.StatusCode != 200 {
		log.Println("request error")
		log.Println(string(ret))
		return
	}
	log.Println("return xml: ")
	log.Println(msg.OriData)
	if exe.analysis != nil {
		msg.Data = exe.analysis.AnalysisResult(msg.OriData, msg.TransCode)
	}
}

func prepareSOAPMsg(msg *message.Message, wsExe WSExecutor, wsConn *connector.WebserviceConnector) string {
	var wb bytes.Buffer
	bodyXml := wsExe.BaseExecutor.template.GetTemplateFromMessage(*msg)
	wb.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	wb.WriteString("<soap:Envelope xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" ")
	wb.WriteString("xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" ")
	wb.WriteString("xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">")
	wb.WriteString("<soap:Body>")
	if wsConn.MethodNamespace() == "" {
		wb.WriteString("<" + wsConn.MethodName())
		if wsConn.Namespace() != "" {
			wb.WriteString(" xmlns=\"" + wsConn.Namespace() + "\"")
		}
		wb.WriteString(">")
	} else {
		wb.WriteString("<methodNs:" + wsConn.MethodName())
		wb.WriteString(" xmlns:methodNs=\"" + wsConn.MethodNamespace() + "\">")
	}

	// 设置了参数，那么发送参数列表
	if wsConn.ParamNames() != "" {
		paramArray := strings.Split(wsConn.ParamNames(), ",")
		for _, paramName := range paramArray {
			wb.WriteString("<" + paramName + ">")
			wb.WriteString("<![CDATA[" + msg.Params[paramName] + "]]>")
			wb.WriteString("</" + paramName + ">")
		}
	} else { // 未设置参数，即tpl内设置报文内容，程序不再处理相关动作
		wb.WriteString(bodyXml)
	}

	if wsConn.MethodNamespace() == "" {
		wb.WriteString("</" + wsConn.MethodName() + ">")
	} else {
		wb.WriteString("</methodNs:" + wsConn.MethodName() + ">")
	}
	wb.WriteString("</soap:Body>")
	wb.WriteString("</soap:Envelope>")
	return wb.String()
}
