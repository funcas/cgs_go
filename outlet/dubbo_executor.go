// dubbo executor
package outlet

import "github.com/funcas/cgs/message"

var _ Executor = (*DubboExecutor)(nil)

type DubboExecutor struct {
	BaseExecutor
}

func (exe DubboExecutor) CanExecute(transCode string) bool {
	return true
}

func (exe DubboExecutor) Execute(msg *message.Message) {
	return
}
