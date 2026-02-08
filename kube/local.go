package kube

import (
  "fmt"

  "github.com/yoophox/emei/cfg"
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/kube/errors"
)

func lookupServerInLocal(svc string) (*Server, error) {
  net, err := lookupNetInLocal(svc)
  if err != nil {
    return nil, err
  }

  var ip []string
  dSvcName := getDeployedSeviceName(svc)
  if dSvcName == "" {
    return nil, errs.ErrorF(errs.ErrId(errors.ERR_ID_KUBE_EMPTY_SERVICE_NAME))
  }
  if _hosts == nil {
    err := initHost()
    if err != nil {
      return nil, err
    }
  }
  err = _hosts.Scan("ips."+dSvcName, &ip)
  if err != nil {
    return nil, err
  }
  if len(ip) == 0 {
    return nil, fmt.Errorf("no local ip for: %s", svc)
  }

  return &Server{Net: net, IP: ip[0]}, nil
}

func lookupNetInLocal(svc string) (*Net, error) {
  svcCfg, err := getSvcCfg(svc, "local")
  if err != nil {
    return nil, err
  }

  return lookupNetFromCfg(svc, svcCfg)
}

func lookupIPInLocal(svc string) (string, error) {
  return "", nil
}

func lookupEPTsInLocal(svc string) (ips []string, err error) {
  return nil, nil
}

var _hosts cfg.Config

// initHost ...
func initHost() error {
  uri := cfg.BuildCfgURI(cfg.CFG_SOURCE_LOCAL, *_localDir+"/host.json")
  var err error
  _hosts, err = cfg.New(uri)
  if err != nil {
    return err
  }

  return nil
}
