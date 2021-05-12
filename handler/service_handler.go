package handler

import (
	"context"
	"errors"
	"log"

	"github.com/funcas/cgs/message"

	"github.com/funcas/cgs/container"

	"github.com/funcas/cgs/gen-go/process"
)

type EntryServiceServiceHandler struct {
}

func NewEntryServiceHandler() *EntryServiceServiceHandler {
	return &EntryServiceServiceHandler{}
}

func (h EntryServiceServiceHandler) Execute(ctx context.Context, transCode string) (_r *process.Resp, _err error) {
	log.Println("invoke " + transCode)
	if len(transCode) <= 0 {
		_err = errors.New("transCode is required")
		return
	}
	transMap, err := container.App().SafeGet(container.OUTLET_NAME)
	if err != nil {
		_err = err
		return
	}
	t := transMap.(container.TransCodeMap)
	outlet, e := t.GetOutlet(transCode)
	if e != nil {
		_err = e
		return
	}
	msg := &message.Message{
		TransCode: transCode,
	}
	outlet.Executor().Execute(msg)
	_r = &process.Resp{TransCode: transCode, Data: msg.OriData}
	return
}

func (h EntryServiceServiceHandler) ExecuteWithParams(ctx context.Context, transCode string, params map[string]string) (_r *process.Resp, _err error) {
	return &process.Resp{}, nil
}

func (h EntryServiceServiceHandler) Reload(ctx context.Context) (_r *process.Resp, _err error) {
	return &process.Resp{}, nil
}

//func GetBytes(key interface{}) ([]byte, error) {
//	var buf bytes.Buffer
//	enc := gob.NewEncoder(&buf)
//	err := enc.Encode(key)
//	if err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}
