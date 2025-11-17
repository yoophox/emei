package etcintra

var (
  EtcBcks = map[string]func([]string, ...*Option) Client{}
  EtcCli  Client
)
