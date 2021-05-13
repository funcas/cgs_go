package container

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/funcas/cgs/manager"

	"github.com/funcas/cgs/tpl"

	"github.com/funcas/cgs/outlet"

	"github.com/funcas/cgs/connector"

	"github.com/valyala/fastjson"

	"github.com/sarulabs/di/v2"
)

const (
	OutletName   = "outlets"
	DispatchName = "dispatch"
)

func Build() {
	builder, err := di.NewBuilder()
	if err != nil {
		log.Fatal(err.Error())
	}
	builder.Set(tpl.DefaultTemplateServiceName, tpl.NewDefaultTemplateService("./conf/template"))
	registerConnectors(builder)
	registerOutlets(builder)
	registerDispatch(builder)

	manager.InitDiFactory(builder.Build())

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
		tplService := string(exe.GetStringBytes("templateService"))
		switch outlet.ExecType(t) {
		case outlet.HttpExec:
			builder.Add(di.Def{
				Name: name,
				Build: func(ctn di.Container) (interface{}, error) {
					c := ctn.Get(conn).(connector.Connector)
					if c.Enabled() {
						return outlet.NewHttpExecutor(c, ctn.Get(tplService).(tpl.TemplateService)), nil
					}
					return nil, errors.New("connector is disabled")
				},
			})
		}
	}
	transMap := make(outlet.TransCodeMap)
	for _, out := range outletSet[1] {
		var acceptTransCodes []string
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
		builder.Set(OutletName, transMap)
	}

}

func registerDispatch(builder *di.Builder) {
	builder.Add(di.Def{
		Name: DispatchName,
		Build: func(ctn di.Container) (interface{}, error) {
			return outlet.NewDispatch(ctn.Get(OutletName).(outlet.TransCodeMap)), nil
		},
	})
}
