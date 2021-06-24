package message

type Message struct {
	TransCode string
	Params    map[string]string
	OriData   string
	Data      string
	ErrMsg    string
}
