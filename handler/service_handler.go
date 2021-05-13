package handler

import (
	"context"
	"errors"
	"log"

	"github.com/funcas/cgs/container"
	"github.com/funcas/cgs/gen-go/process"
	"github.com/funcas/cgs/message"
)

type EntryServiceServiceHandler struct {
}

func NewEntryServiceHandler() *EntryServiceServiceHandler {
	return &EntryServiceServiceHandler{}
}

func (h EntryServiceServiceHandler) Execute(ctx context.Context, transCode string) (_r *process.Resp, _err error) {
	return h.ExecuteWithParams(ctx, transCode, nil)
}

func (h EntryServiceServiceHandler) ExecuteWithParams(ctx context.Context, transCode string, params map[string]string) (_r *process.Resp, _err error) {
	log.Println("invoke " + transCode)
	if len(transCode) <= 0 {
		_err = errors.New("transCode is required")
		return
	}
	msg := &message.Message{
		TransCode: transCode,
		Params:    params,
	}
	dispatch := container.App().Get(container.DispatchName).(*container.Dispatch)
	err := dispatch.Send(msg)
	if err != nil {
		_err = err
	}
	_r = &process.Resp{TransCode: transCode, Data: msg.OriData}
	return
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
