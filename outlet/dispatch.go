package outlet

import (
	"log"

	"github.com/funcas/cgs/message"
)

type Dispatch struct {
	transMap TransCodeMap
}

func NewDispatch(transMap TransCodeMap) *Dispatch {
	return &Dispatch{
		transMap,
	}
}

func (d Dispatch) Send(msg *message.Message) {
	transCode := msg.TransCode
	outlet, e := d.transMap.GetOutlet(transCode)
	if e != nil {
		log.Println(e.Error())
		return
	}
	outlet.Executor().Execute(msg)
	return
}
