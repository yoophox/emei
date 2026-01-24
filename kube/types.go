package kube

import (
  "sync"

  "github.com/yoophox/emei/cfg"
)

type Server struct {
  IP   string
  Net  string
  Port string
}

type Net struct {
  Port string
  Net  string
}

type svcCache struct {
  m    sync.RWMutex
  svcs map[string]cfg.Config
}
