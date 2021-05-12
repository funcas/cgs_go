package container

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/funcas/cgs/outlet"

	"github.com/funcas/cgs/connector"

	"github.com/valyala/fastjson"

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
	registerConnectors(builder)
	registerOutlets(builder)
	app = builder.Build()

}

func Destroy() {
	app.Delete()
}

func readConfig(path string, keys ...string) [][]*fastjson.Value {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	var p fastjson.Parser
	v, e := p.ParseBytes(data)
	if e != nil {
		log.Fatal(e.Error())
	}
	res := make([][]*fastjson.Value, len(keys))
	for i, key := range keys {
		res[i] = v.GetArray(key)
	}
	return res
}

//
//func readOutletConf()

// register connector instances to di container
func registerConnectors(builder *di.Builder) {
	connList := readConfig("./conf/connector.json", "connectors")
	for _, conn := range connList[0] {
		var err error
		t := string(conn.GetStringBytes("type"))
		name := string(conn.GetStringBytes("name"))
		switch connector.ConnType(t) {
		case connector.Http:
			err = builder.Set(name,
				connector.NewHttpConnector(conn))
		case connector.Socket:
			err = errors.New("not support yet")
		case connector.WebService:
			err = errors.New("not support yet")
		default:
			err = builder.Set(name, connector.NewHttpConnector(conn))
		}
		if err != nil {
			log.Fatal(err.Error())
		}

	}
}

func registerOutlets(builder *di.Builder) {
	outletSet := readConfig("./conf/outlet.json", "executors", "outlets")
	for _, exe := range outletSet[0] {
		t := string(exe.GetStringBytes("type"))
		name := string(exe.GetStringBytes("name"))
		conn := string(exe.GetStringBytes("connector"))
		switch outlet.ExecType(t) {
		case outlet.HttpExec:
			builder.Add(di.Def{
				Name: name,
				Build: func(ctn di.Container) (interface{}, error) {
					return outlet.NewHttpExecutor(ctn.Get(conn).(connector.Connector)), nil
				},
			})
		}
	}
	transMap := make(TransCodeMap)
	for _, out := range outletSet[1] {
		acceptTransCodes := []string{}
		name := string(out.GetStringBytes("name"))
		exec := string(out.GetStringBytes("executor"))
		transCodeArr := out.GetArray("acceptTransCodes")
		for _, v := range transCodeArr {
			transCode := string(v.GetStringBytes())
			acceptTransCodes = append(acceptTransCodes, transCode)
			transMap[transCode] = name
		}
		builder.Add(di.Def{
			Name: name,
			Build: func(ctn di.Container) (interface{}, error) {
				return outlet.NewOutlet(name, ctn.Get(exec).(outlet.Executor), acceptTransCodes), nil
			},
		})
	}
	if len(transMap) > 0 {
		builder.Set(OUTLET_NAME, transMap)
	}

}

const OUTLET_NAME = "outlets"

type TransCodeMap map[string]string

func (t TransCodeMap) Exists(transCode string) bool {
	_, exists := t[transCode]
	return exists
}

func (t TransCodeMap) GetOutlet(transCode string) (*outlet.Outlet, error) {
	if !t.Exists(transCode) {
		return nil, errors.New("invalid transCode")
	}

	res, err := app.SafeGet(t[transCode])
	if err != nil {
		return nil, err
	}
	return res.(*outlet.Outlet), nil
}
