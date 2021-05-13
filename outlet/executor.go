package outlet

import (
	"github.com/funcas/cgs/connector"
	"github.com/funcas/cgs/message"
	"github.com/funcas/cgs/tpl"
)

type Executor interface {
	CanExecute(transCode string) bool

	Execute(msg *message.Message)
}

type BaseExecutor struct {
	connector connector.Connector
	template  tpl.TemplateService
}

type ExecType string

const (
	HttpExec ExecType = "HTTP_EXECUTOR"
)
