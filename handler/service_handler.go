package handler

import (
	"context"
	"errors"
	"log"

	"github.com/funcas/cgs/outlet"

	"github.com/funcas/cgs/manager"

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
	dispatch := manager.GetDiFactory().Get(container.DispatchName).(*outlet.Dispatch)
	dispatch.Send(msg)
	_r = &process.Resp{TransCode: transCode, Data: msg.OriData, ErrorMsg: msg.RetMsg}
	return
}

func (h EntryServiceServiceHandler) Reload(ctx context.Context) (_r *process.Resp, _err error) {
	return &process.Resp{}, nil
}
