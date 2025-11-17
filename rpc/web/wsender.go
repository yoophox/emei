package web

import (
  "fmt"
  "net/http"
  "reflect"
  "time"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/kube/resolver"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/call"
  "github.com/yolksys/emei/rpc/coder"
  "github.com/yolksys/emei/rpc/rpcabs"
)

type webres struct {
  r *http.Response
}

func (w *webres) RValues(rTyp ...reflect.Type) ([]reflect.Value, error) {
  dec := newDec(w.r.Body)
  defer dec.Release()

  h, err := coder.Header(dec)
  if err != nil {
    return coder.DefaultValues(nil, rTyp...), fmt.Errorf("fail:get header, par:{%s}", err.Error())
  }

  if h.Code != 0 {
    return coder.DefaultValues(nil, rTyp...), fmt.Errorf("fail:h.code,reason:%s", h.Reason)
  }

  res, err := coder.Values(dec, rTyp...)
  if err == nil {
    return res, nil
  }

  return coder.DefaultValues(nil, rTyp...), fmt.Errorf("fail:get values, reason:%s", err.Error())
}

func (w *webres) Close() {
}

type websender struct {
  http http.Client
  rpcabs.SndImp
}

func (w *websender) Send(e env.Env, svc, met string, args ...any) call.Response {
  defer e.Return()

  svc_, err := w.GetSvc(svc)
  e.AssertErr(err)

  enc := newEnc()
  defer enc.Release()
  h := e.GetMsgHeader()
  err = enc.Encode(h)
  e.AssertErr(err)
  for _, value := range args {
    err = enc.Encode(value)
    e.AssertErr(err)
  }

  r, err := w.http.Post(fmt.Sprintf("%s://%s:%s/%s", svc_.Trans, svc_.Ip, svc_.Port,
    met), "application/json", enc.getReader())
  e.AssertErr(err, func() { w.Delete(svc) })
  return &webres{r}
}

// newWebSender ...
func newWebSender() *websender {
  httpTransport := &http.Transport{
    Proxy:                 http.ProxyFromEnvironment,
    ForceAttemptHTTP2:     true,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
    TLSClientConfig:       pki.NewClientTlsConfig(),
    DialContext:           resolver.DialContent,
  }
  return &websender{
    http: http.Client{
      Transport: httpTransport,
      Timeout:   10 * time.Second,
    },

    SndImp: rpcabs.SndImp{
      For: "web",
    },
  }
}
