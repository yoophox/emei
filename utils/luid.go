package utils

import (
  "fmt"

  "github.com/sony/sonyflake/v2"
)

// LUID ...
func LUID() int64 {
  l, err := _luidgener.NextID()
  AssertErr(err)
  return l
}

// LUIDHex ...
func LUIDHex() string {
  return fmt.Sprintf("%016X", LUID())
}

var _luidgener *sonyflake.Sonyflake

func init() {
  var err error
  _luidgener, err = sonyflake.New(sonyflake.Settings{})
  AssertErr(err)
}
