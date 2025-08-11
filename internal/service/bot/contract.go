package bot

type CommandFactory interface {
	GetCommand(name string) Command
}

type Command interface {
	Execute(cmd string) (string, error)
	SetContext(ctx interface{})
}

type Metrics interface {
	Count(name string, value int64)
	Disconnect()
	Duration(timestamp int64, name string)
	Gauge(name string, value float64)
	Increment(name string)
}
