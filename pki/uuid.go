package pki

import (
  "github.com/sony/sonyflake/v2"
  "github.com/yoophox/emei/utils"
)

// github.com/sony/sonyflake/v2

// LocalUid ...
func LocalUid() int64 {
  i, err := _localIdGener.NextID()
  if err != nil {
    utils.AssertErr(err)
  }
  return i
}

// UUID ...
func UUID() int64 {
  i := <-_uuidCh
  return i
}

// UUIDS ...
// @param num: <= 100
func UUIDS(num int) []int64 {
  return nil
}

func init() {
  var err error
  _localIdGener, err = sonyflake.New(
    sonyflake.Settings{
      MachineID: func() (int, error) {
        return 1000, nil
      },
    },
  )
  if err != nil {
    utils.AssertErr(err)
  }
}

var (
  _localIdGener *sonyflake.Sonyflake
  _uuidCh       = make(chan int64, 1000)
)
