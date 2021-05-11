package executor

import "github.com/funcas/cgs/message"

type Executor interface {
	CanExecute(transCode string) bool

	Execute(msg message.Message)
}
