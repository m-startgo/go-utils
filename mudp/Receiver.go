package mudp

type Receiver struct {
	Addr      string
	MultiCore bool
}

func NewReceiver() *Receiver {
	return &Receiver{}
}

func (c *Receiver) Start() {
}
