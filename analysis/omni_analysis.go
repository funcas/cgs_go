package analysis

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/funcas/cgs/manager"

	"github.com/jf-tech/omniparser/transformctx"

	"github.com/jf-tech/omniparser"
)

const OmniAnalysisName = "OmniAnalyser"

type OmniAnalysis struct {
}

func NewOmniAnalysis() *OmniAnalysis {
	return &OmniAnalysis{}
}

func (o OmniAnalysis) AnalysisResult(result, transCode string) string {
	schema, e := getOmniSchema(transCode)
	// sth err occur
	if e != nil {
		log.Println(e.Error())
		return ""
	}
	transform, err := schema.NewTransform(transCode, strings.NewReader(result), &transformctx.Ctx{})
	if err != nil {
		log.Println(err.Error())
	}
	var ret []byte
	for {
		output, err := transform.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err.Error())
		}
		ret = append(ret, output...)
	}
	return string(ret)
}

func getOmniSchema(transCode string) (schema omniparser.Schema, e error) {
	schemaFile := manager.TransFormMap()[transCode]
	if schemaFile == nil {
		errorMsg := fmt.Sprintf("could not find transform schema file %s.json", transCode)
		e = errors.New(errorMsg)
		return
	}
	schema, err := omniparser.NewSchema(transCode, schemaFile)
	if err != nil {
		e = err
		return
	}
	return
}
