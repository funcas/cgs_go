package container

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/funcas/cgs/connector"
	"github.com/funcas/cgs/model"

	"github.com/sarulabs/di/v2"
)

var app di.Container

func App() di.Container {
	return app
}

func Build() {
	builder, err := di.NewBuilder()
	if err != nil {
		log.Fatal(err.Error())
	}

	connList := readConnectors()
	for _, conn := range connList.Connectors {
		err := builder.Set(conn.Name, connector.NewHttpConnector(conn))
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	//builder.Add(di.Def{
	//	Name: "conn3",
	//	Build: func(ctn di.Container) (interface{}, error) {
	//		return connector.NewHttpConnector("conn3", "http://www.google.com"), nil
	//	},
	//	Close: func(obj interface{}) error {
	//		fmt.Printf(">>> ===== %s is being destroied ===== <<<\n", obj.(*connector.HttpConnector).Name())
	//		return nil
	//	},
	//})
	app = builder.Build()

}

func Destroy() {
	app.Delete()
}

type ConnectorList struct {
	Connectors []model.HttpConnectorVO
}

func readConnectors() *ConnectorList {
	connList := &ConnectorList{}
	data, err := ioutil.ReadFile("./conf/connector.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	json.Unmarshal(data, connList)

	return connList
}
