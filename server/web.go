package svr

import (
  "net/http"
  "reflect"
  "strings"

  "github.com/gorilla/websocket"
  "github.com/yoophox/emei/env"
)

// serveHttp ...
func serveHttp(w http.ResponseWriter, r *http.Request) {
  topic := strings.TrimPrefix(r.URL.Path, "/")
  pos := strings.Index(topic, ".")
  if pos <= 0 {
    w.Write([]byte("path have no dot:" + topic))
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  var tja env.Tjatse
  midc, err := r.Cookie(RPC_TID_COOKIES_NAME)
  if err == nil {
    tja.Mid = midc.Value
  }
  jwtc, err := r.Cookie(RPC_JWT_COOKIES_NAME)
  if err == nil {
    tja.Jwt = jwtc.Value
  }

  e := env.New(&tja)
  var wsp *webResponsImpl = nil
  defer func() {
    err := e.Err()
    e.Release()
    if err != nil {
      w.Write([]byte(err.Error()))
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    if wsp == nil {
      return
    }

    for k, v := range wsp.headers {
      w.Header().Add(k, v)
    }
    for _, v := range wsp.cookies {
      http.SetCookie(w, v)
    }
    for _, v := range wsp.content {
      w.Write(v)
    }
    w.WriteHeader(http.StatusOK)
  }()
  defer e.Trace(topic)

  rcvr, ok := _webRecvs[topic[:pos]]
  e.AssertBool(ok, ERR_ID_RPC_CALLINFO_RCVR, "no rcvr for topic: %s", topic)
  m, ok := rcvr.funcs[topic[pos+1:]]
  e.AssertBool(ok, ERR_ID_RPC_CALLINFO_RCVR, "no rcvr for func: %s", topic)
  ptyps := rcvr.params[topic[pos+1:]]
  params := []reflect.Value{rcvr.value, reflect.ValueOf(e)}
  if len(ptyps) == 3 {
    e.AssertBool(r.Header.Get("upgrade") == "websocket", ERR_ID_WEB_NOT_WS_CALL, "met is for ws call but not ws request: %s", topic)
    ws, err := newwebSocket(w, r)
    e.AssertErr(err, ERR_ID_WEB_NEW_SOCKET)
    params = append(params, reflect.ValueOf(ws))
  } else {
    wsp = newWebResponse()
    params = append(params, reflect.ValueOf(wsp), reflect.ValueOf(r))
  }

  m.Call(params)

  // cntTyp := r.Header.Get("content-type")
  // if r.Header.Get("upgrade") == "websocket" {
  // } else if cntTyp == "multipart/form-data" {
  //   err := r.ParseMultipartForm(1000 * 1000 * 100)
  //   if err != nil {
  //     w.Write([]byte(err.Error()))
  //     w.WriteHeader(http.StatusBadRequest)
  //     return
  //   }
  //   // stream = (*upFiler)(r.MultipartForm)
  // } else if cntTyp == "application/json" {
  // } else {
  //   w.Write([]byte("wrong content type:" + cntTyp))
  //   w.WriteHeader(http.StatusBadRequest)
  // }
}

// newwebSocket ...
func newwebSocket(w http.ResponseWriter, r *http.Request) (WebSock, error) {
  up := websocket.Upgrader{}
  return up.Upgrade(w, r, w.Header())
}
