package container

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/funcas/cgs/analysis"

	"github.com/funcas/cgs/manager"

	"github.com/funcas/cgs/tpl"

	"github.com/funcas/cgs/outlet"

	"github.com/funcas/cgs/connector"

	"github.com/valyala/fastjson"

	"github.com/sarulabs/di/v2"
)

const (
	OutletName    = "outlets"
	DispatchName  = "dispatch"
	TemplateDir   = "./conf/template"
	ConnectorConf = "./conf/connector.json"
)

func Build() {
	builder, err := di.NewBuilder()
	if err != nil {
		log.Fatal(err.Error())
	}
	// register connectors to di container
	registerConnectors(builder)
	// register templateService to di container
	registerTmplService(builder)
	// register analysisService to di container
	registerAnalysisService(builder)
	// register outlets to di container
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

func registerTmplService(builder *di.Builder) {
	builder.Set(tpl.DefaultTemplateServiceName, tpl.NewDefaultTemplateService(TemplateDir))
}

// register connector instances to di container
func registerConnectors(builder *di.Builder) {
	connList := readConfig(ConnectorConf, "connectors")
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
			err = builder.Set(name, connector.NewWebserviceConnector(conn))
		default:
			err = builder.Set(name, connector.NewHttpConnector(conn))
		}
		if err != nil {
			log.Fatal(err.Error())
		}

	}
}

//
func registerAnalysisService(builder *di.Builder) {
	builder.Set(analysis.OmniAnalysisName, analysis.NewOmniAnalysis())
	builder.Set(analysis.DefaultAnalysisName, analysis.NewDefaultAnalysis())
}

func registerOutlets(builder *di.Builder) {
	outletSet := readConfig("./conf/outlet.json", "executors", "outlets")
	for _, exe := range outletSet[0] {
		t := string(exe.GetStringBytes("type"))
		name := string(exe.GetStringBytes("name"))
		conn := string(exe.GetStringBytes("connector"))
		analy := string(exe.GetStringBytes("analysisService"))
		tplService := string(exe.GetStringBytes("templateService"))
		switch outlet.ExecType(t) {
		case outlet.HttpExec:
			builder.Add(di.Def{
				Name: name,
				Build: func(ctn di.Container) (interface{}, error) {
					connector := ctn.Get(conn).(connector.Connector)
					var analyser = ctn.Get(analysis.DefaultAnalysisName).(analysis.Analyser)
					var templater = ctn.Get(tpl.DefaultTemplateServiceName).(tpl.TemplateService)
					if analy != "" {
						analyser = ctn.Get(analy).(analysis.Analyser)
					}
					if tplService != "" {
						templater = ctn.Get(tplService).(tpl.TemplateService)
					}

					if connector.Enabled() {
						return outlet.NewHttpExecutor(connector,
							templater,
							analyser), nil
					}
					return nil, errors.New("connector is disabled")
				},
			})
		case outlet.WSExec:
			builder.Add(di.Def{
				Name: name,
				Build: func(ctn di.Container) (interface{}, error) {
					connector := ctn.Get(conn).(connector.Connector)
					var analyser = ctn.Get(analysis.DefaultAnalysisName).(analysis.Analyser)
					var templater = ctn.Get(tpl.DefaultTemplateServiceName).(tpl.TemplateService)
					if analy != "" {
						analyser = ctn.Get(analy).(analysis.Analyser)
					}
					if tplService != "" {
						templater = ctn.Get(tplService).(tpl.TemplateService)
					}

					if connector.Enabled() {
						return outlet.NewWSExecutor(connector,
							templater,
							analyser), nil
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
