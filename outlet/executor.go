package outlet

import (
	"github.com/funcas/cgs/connector"
	"github.com/funcas/cgs/message"
)

type Executor interface {
	CanExecute(transCode string) bool

	Execute(msg *message.Message)
}

type BaseExecutor struct {
	connector connector.Connector
}

type ExecType string

const (
	HttpExec ExecType = "HTTP_EXECUTOR"
)
