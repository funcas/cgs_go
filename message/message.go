package message

type Message struct {
	TransCode string
	Params    map[string]string
	OriData   string
	RetMsg    string
	RetCode   string
}
