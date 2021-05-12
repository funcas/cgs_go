package model

type Executors struct {
	Name            string `json:"name"`
	Connector       string `json:"connector"`
	AnalysisService string `json:"analysisService"`
	TemplateService string `json:"templateService"`
}

type Outlets struct {
	Name             string   `json:"name"`
	Executors        []string `json:"executors"`
	AcceptTransCodes []string `json:"acceptTransCodes"`
}

type OutletConf struct {
	Executors []Executors `json:"executors"`
	Outlets   []Outlets   `json:"outlets"`
}
