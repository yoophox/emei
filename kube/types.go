package kube

import (
  "sync"

  "github.com/yoophox/emei/cfg"
)

type Server struct {
  IP string
  *Net
}

type Net struct {
  Ports map[string]*ServicePort

  ports []*ServicePort
  svc   string
}

type ServicePort struct {
  Name, Port, TargetPort, NodePort string
}

type svcCache struct {
  m    sync.RWMutex
  svcs map[string]cfg.Config
}
