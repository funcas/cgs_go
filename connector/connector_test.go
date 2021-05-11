package connector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/funcas/cgs/model"
)

type ConnectorList struct {
	Connectors []model.HttpConnectorVO
}

func TestReadConnectorFromJson(t *testing.T) {
	conns := &ConnectorList{}
	data, err := ioutil.ReadFile("../conf/connector.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	json.Unmarshal(data, conns)
	conn := NewHttpConnector(conns.Connectors[0])
	t.Log(conn.Timeout())
	t.Log(conn)
}
