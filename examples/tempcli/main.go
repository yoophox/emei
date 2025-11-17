package main

import (
  "errors"
  "fmt"
  "time"

  "github.com/yolksys/emei"
  "yolk.com/pdt/temp"
)

func main() {
  defer time.Sleep(time.Second * 13)
  emei.Log("kkkkk", "FAFDFGKNJLLIJ")
  req := temp.Request{Name: "temp client"}
  e := emei.NewEnv("tempclient", "req", nil, nil)
  defer e.Finish()

  r, err := emei.Call1[*temp.Response](e, "temp", "server.hello", req)
  fmt.Println(r, err)

  r, err = emei.Call1[*temp.Response](e, "temp", "server.hello", req)
  fmt.Println(r, err)

  r, err = emei.Call1[*temp.Response](e, "temp", "server.hello", req)
  fmt.Println(r, err)

  // io_, err := emei.CallWithWStream(e, "temp", "server.stream", &req)
  // if err != nil {
  //   fmt.Println("error", err)
  //   return
  // }
  // fmt.Println("UUUUUUUUUUUUUUU success")

  // for i := 0; i < 10; i++ {
  //   n, err := io_.Write([]byte(fmt.Sprintf("afafsadfasdfasfa: %d", i)))
  //   if err != nil {
  //     fmt.Println("UUUUUUUUUUUUUU err", err)
  //     break
  //   } else {
  //     fmt.Println("UUUUUUUUUUUUU len", n)
  //   }
  // }
  // err = io_.Close()
  // fmt.Println("UUUUUUUUUUUUUUUUUUU close", err)

  // // time.Sleep(time.Second)
  // io_, err = emei.CallWithWStream(e, "temp", "server.stream", &req)
  // if err != nil {
  //   fmt.Println("error", err)
  //   return
  // }

  // fmt.Println("UUUUUUUUUUUUUUU success")

  // for i := 0; i < 10; i++ {
  //   n, err := io_.Write([]byte(fmt.Sprintf("afafsadfasdfasfa: %d", i)))
  //   if err != nil {
  //     fmt.Println("UUUUUUUUUUUUUU err", err)
  //     break
  //   } else {
  //     fmt.Println("UUUUUUUUUUUUU len", n)
  //   }
  // }
  // err = io_.Close()
  // fmt.Println("UUUUUUUUUUUUUUUUUUU close", err)

  // // time.Sleep(time.Second)
  // io_, err = emei.CallWithWStream(e, "temp", "server.stream", &req)
  // if err != nil {
  //   fmt.Println("error", err)
  //   return
  // }

  // fmt.Println("UUUUUUUUUUUUUUU success")

  // for i := 0; i < 10; i++ {
  //   n, err := io_.Write([]byte(fmt.Sprintf("afafsadfasdfasfa: %d", i)))
  //   if err != nil {
  //     fmt.Println("UUUUUUUUUUUUUU err", err)
  //     break
  //   } else {
  //     fmt.Println("UUUUUUUUUUUUU len", n)
  //   }
  // }
  // err = io_.Close()
  // fmt.Println("UUUUUUUUUUUUUUUUUUU close", err)

  // fmt.Println(name())
}

// name ...
func name() (err error) {
  defer func() {
    err = errors.New("sdfasdfasdfasdafdafs")
  }()

  return nil
}
