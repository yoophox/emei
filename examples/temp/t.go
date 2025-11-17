package main

import (
  "fmt"
  "reflect"
)

// ttt ...
func ttt() {
  type data struct {
    firstName string
    lastName  string
  }
  var v []data
  var v_ any = &v
  vvv := reflect.ValueOf(v_)
  v__ := reflect.TypeOf(v_)
  t := reflect.MakeSlice(v__.Elem(), 1, 4)
  t.Index(0).Field(0).SetString("asdfasdfasdffasfd")
  vvv.Elem().Set(t)

  fmt.Println(v)
}
