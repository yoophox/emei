package web

import (
  "errors"
  "fmt"
  "net/http"
  "reflect"
  "strings"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/rpc/coder"
)

// route ...
func (s *webrpc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // buf := utils.NewBuffer()
  // io.Copy(buf, r.Body)b
  // r.Body.Close()
  cntTyp := r.Header.Get("content-type")
  fmt.Printf("UUUUUUUUUUUUUUUUUUU:%s------------%s----------%s\n", cntTyp, r.Header.Get("accept"), r.Header.Get("upgrade"))
  var stream any

  p := strings.Trim(r.URL.Path, "/")
  pos := strings.Index(p, "/")
  metPath := ""
  if pos <= 0 {
    metPath = p
  } else {
    metPath = p[:pos]
  }

  e := env.New("web", metPath, newEnc(w), newDec(r))
  defer e.Finish()
  e.Assert()
  wh_ := func() { w.WriteHeader(http.StatusInternalServerError) }

  if r.Header.Get("upgrade") == "websocket" {
    c, err := newWebSock(w, r)
    e.AssertErr(err, wh_)
    stream = c
  } else if cntTyp == "multipart/form-data" {
    err := r.ParseMultipartForm(1024 * 1024 * 10)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      sss := fmt.Sprintf("fail:http.ParseMultipartForm, reason:%s, path:%s", err.Error(), r.URL.Path)
      w.Write([]byte(sss))
      log.Error("msg", sss)
      return
    }
    stream = (*upFiler)(r.MultipartForm)
  } else if pos > 0 {
    stream = newDnFiler(w, p[pos+1:])
  } else {
    stream = nil
  }

  _path := strings.Split(metPath, ".")
  if len(_path) != 2 {
    e.AssertErr(errors.New("fail: path, path:"+r.URL.Path), wh_)
  }

  recv, ok := s.Recvs[_path[0]]
  if !ok {
    e.AssertErr(errors.New("fail: no recver, name:"+_path[0]), wh_)
  }

  met, ok := recv.Mets[_path[1]]
  if !ok {
    e.AssertErr(errors.New("fail: no method, name:"+_path[1]), wh_)
  }
  parms := coder.ParseParam(e, met, recv)
  e.Assert()
  e.PrintParams(parms...)
  if stream != nil {
    parms = append(parms, reflect.ValueOf(stream))
  }

  f := met.Method.Func
  // res := f.Call(parms)
  e.SetReV(f.Call(parms))
}
