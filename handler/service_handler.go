package handler

import (
	"context"
	"log"

	"github.com/funcas/cgs/executor"
	"github.com/funcas/cgs/gen-go/process"
	"github.com/funcas/cgs/message"

	"github.com/funcas/cgs/connector"
	"github.com/funcas/cgs/container"
)

type EntryServiceServiceHandler struct {
}

func NewEntryServiceHandler() *EntryServiceServiceHandler {
	return &EntryServiceServiceHandler{}
}

func (h EntryServiceServiceHandler) Execute(ctx context.Context, transCode string) (_r *process.Resp, _err error) {
	log.Println("invoke " + transCode)

	hc, _ := container.App().SafeGet(transCode)

	exec := &executor.HttpExecutor{}
	exec.SetHttpConnector(hc.(*connector.HttpConnector))
	exec.Execute(message.Message{TransCode: ""})
	return &process.Resp{TransCode: transCode, Data: map[string]string{"hello": transCode}}, nil
}

func (h EntryServiceServiceHandler) ExecuteWithParams(ctx context.Context, transCode string, params process.Data) (_r *process.Resp, _err error) {
	return &process.Resp{}, nil
}
