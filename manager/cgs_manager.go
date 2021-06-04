package manager

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var transformMap = make(map[string]*os.File)

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
		transformMap[fileName[0]] = schemaFile
	}
}

func TransFormMap() map[string]*os.File {
	return transformMap
}
