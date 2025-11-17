package web

import (
  "context"
  "errors"
  "net/http"

  "github.com/yolksys/emei/log"
)

func (s *webrpc) Name() string {
  return s.RpcImp.Name
}

func (s *webrpc) Start() error {
  //_https._wg.Add(1)
  s.real = &http.Server{
    Addr:    ":" + s.Port,
    Handler: s,
  }

  go func() {
    defer func() {
      //_https._wg.Done()
    }()

    var err error
    if s.Trans == "https" {
      log.Event("web https listen port", s.Port)
      err = s.real.ListenAndServeTLS(s.Cert, s.Key)
    } else {
      log.Event("web http listen port", s.Port)
      err = s.real.ListenAndServe()
    }

    if err != nil {
      s.Err <- errors.New("status:stoped, msg:" + err.Error())
    }
  }()

  return nil
}

func (s *webrpc) Close() {
  s.IsClosed_ = true
  go s.real.Shutdown(context.Background())
}

// _new ...
func NewRpc() *webrpc {
  s := &webrpc{}
  s.RpcImp.Name = "web"
  return s
}
