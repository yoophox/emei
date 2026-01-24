package kube

import (
  "fmt"

  "github.com/yoophox/emei/cfg"
)

func lookupServerInLocal(svc string) (*Server, error) {
  net, err := lookupNetInLocal(svc)
  if err != nil {
    return nil, err
  }

  var ip []string
  dSvcName := getDeployedSeviceName(svc)
  if _hosts == nil {
    initHost()
  }
  err = _hosts.Scan("ips."+dSvcName, &ip)
  if err != nil {
    return nil, err
  }
  if len(ip) == 0 {
    return nil, fmt.Errorf("no local ip for: %s", svc)
  }

  return &Server{Port: net.Port, Net: net.Net, IP: ip[0]}, nil
}

func lookupNetInLocal(svc string) (*Net, error) {
  svcCfg, err := getSvcCfg(svc, "local")
  if err != nil {
    return nil, err
  }

  return lookupNetFromCfg(svcCfg)
}
func lookupIPInLocal(svc string) (string, error)
func lookupEPTsInLocal(svc string) (ips []string, err error) // get all ips of a service

var _hosts cfg.Config

// initHost ...
func initHost() {
  uri := cfg.BuildCfgURI(cfg.CFG_SOURCE_LOCAL, _localDir+"/host.json")
  var err error
  _hosts, err = cfg.New(uri)
  if err != nil {
    panic(err)
  }
}
