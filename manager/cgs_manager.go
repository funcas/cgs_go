package manager

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jf-tech/omniparser"
)

var transformMap = make(map[string]omniparser.Schema)

const transformDir = "/conf/transform/"

func LoadTransformer() {
	pwd, _ := os.Getwd()
	f := pwd + transformDir
	log.Println(pwd)
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(f)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range fileInfoList {
		fileName := strings.Split(v.Name(), ".")
		schemaFile, err := os.Open(f + v.Name())
		if err != nil {
			log.Fatal(err.Error())
		}
		schema, err := omniparser.NewSchema(fileName[0], schemaFile)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		transformMap[fileName[0]] = schema
	}
}

func TransFormMap() map[string]omniparser.Schema {
	return transformMap
}
