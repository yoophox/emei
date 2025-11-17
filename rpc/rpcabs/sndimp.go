package rpcabs

import (
  "fmt"
  "io"
  "sync"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/kube"
  "github.com/yolksys/emei/utils"
)

type SvcInfo struct {
  Ip, Port, Trans string
  Client          any
}

type SndImp struct {
  svcs map[string]*SvcInfo
  mux  sync.Mutex
  For  string // impled rpc eg "web"
}

func (s *SndImp) GetSvc(name string, f ...func(*SvcInfo) error) (*SvcInfo, error) {
  svc, ok := s.svcs[name]
  if ok {
    return svc, nil
  }

  s.mux.Lock()
  defer s.mux.Unlock()
  svc, ok = s.svcs[name]
  if ok {
    return svc, nil
  }
  t, ip, port, err := kube.Lookup(name, s.For)
  if err != nil {
    return nil, fmt.Errorf("fail:getsvc, par:{%s}", err.Error())
  }

  svc = &SvcInfo{
    Ip:    ip,
    Port:  port,
    Trans: t,
  }
  if !utils.IsIpv4(ip) {
    svc.Ip = "[" + ip + "]"
  }
  for _, value := range f {
    err = value(svc)
    if err != nil {
      return nil, fmt.Errorf("fail:getsvc, reason:newclient, par:{%s}, ip:%s, port:%s, trans:%s",
        err.Error(), svc.Ip, svc.Port, svc.Trans)
    }

  }
  if s.svcs == nil {
    s.svcs = make(map[string]*SvcInfo)
  }
  s.svcs[name] = svc
  return svc, nil
}

func (s *SndImp) Delete(name string) {
  s.mux.Lock()
  defer s.mux.Unlock()
  delete(s.svcs, name)
}

func (s *SndImp) SendWithStream(e env.Env, svc, met string, args ...any) io.ReadWriteCloser {
  e.AssertErr(fmt.Errorf("fail:GetRStream, msg:%s don'// TODO: implement", s.For))
  return nil
}
