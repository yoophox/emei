package main

import (
  "fmt"
  "io"

  "github.com/yolksys/emei"
  "yolk.com/pdt/temp"
)

type Server struct{}

var i = 0

func (s *Server) Hello(e emei.Env, req *temp.Request) *temp.Response {
  defer e.Return()
  if i > 0 && i < 2 {
    i++
    e.AssertErr(fmt.Errorf("fail:Hello"))
  }
  i++
  return &temp.Response{Echo: "hello " + req.Name + fmt.Sprintf(" %d", i)}
  // return &temp.Response{Echo: "hello " + req.Name}, nil
}

func (s *Server) Stream(e emei.Env, req *temp.Request, io_ io.Reader) {
  defer e.Return()

  buf := [1000]byte{}
  for {
    fmt.Println("UUUUUUUUUUUUUUUUUUUUUU start read")
    n, err := io_.Read(buf[:])
    if err != nil {
      break
    }

    fmt.Println("UUUUUUUU data", string(buf[:n]))
  }

  return
}

func (s *Server) Web(e emei.Env) *temp.Response {
  defer e.Return()

  return &temp.Response{Echo: "hello web"}
}

var _sssss = &Server{}
