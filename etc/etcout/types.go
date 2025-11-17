package etcout

type Event interface {
  EType() string // new/update/delete
  EKey() string
  EValue() string
  EVersion() int64
}
