package tpl

import (
	"log"

	"github.com/flosch/pongo2/v4"
	"github.com/funcas/cgs/message"
)

type TemplateService interface {
	GetTemplate(transCode string) string

	GetTemplateWithParams(transCode string, params map[string]string) string

	GetTemplateFromMessage(msg message.Message) string
}

const DefaultTemplateServiceName = "DefaultTemplateService"

type DefaultTemplateService struct {
	baseDir string
	tplSet  *pongo2.TemplateSet
}

func (t DefaultTemplateService) GetTemplate(transCode string) string {

	return t.GetTemplateWithParams(transCode, nil)
}

func (t DefaultTemplateService) GetTemplateWithParams(transCode string, params map[string]string) string {
	tpl := pongo2.Must(t.tplSet.FromCache(transCode + ".tpl"))
	ctx := pongo2.Context{}
	for k, v := range params {
		ctx[k] = v
	}
	out, err := tpl.Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func (t DefaultTemplateService) GetTemplateFromMessage(msg message.Message) string {
	return t.GetTemplateWithParams(msg.TransCode, msg.Params)
}

func NewDefaultTemplateService(baseDir string) *DefaultTemplateService {
	tplLoader := pongo2.MustNewLocalFileSystemLoader(baseDir)

	// DefaultSet is a set created for you for convinience reasons.
	tplSet := pongo2.NewSet("tplSet", tplLoader)
	return &DefaultTemplateService{
		baseDir,
		tplSet,
	}

}
