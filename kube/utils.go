package kube

import (
  "fmt"

  "github.com/yoophox/emei/cfg"
)

// getSvcCfg ...
func getSvcCfg(svc, uriType string) (cfg.Config, error) {
  _svcCfgs.m.RLock()
  _c, ok := _svcCfgs.svcs[svc]
  _svcCfgs.m.RUnlock()
  if ok {
    return _c, nil
  }

  _svcCfgs.m.Lock()
  defer _svcCfgs.m.Unlock()

  _c, ok = _svcCfgs.svcs[svc]
  if ok {
    return _c, nil
  }

  if svc == "@@self" {
    return _selfSvcCfg, nil
  }

  dSvcName := getDeployedSeviceName(svc)
  uri := ""
  switch uriType {
  case "local":
    uri = fmt.Sprintf("%s~%s/%s", cfg.CFG_SOURCE_LOCAL,
      _localDir+"/"+dSvcName+".yaml")
  case "etc":
    uri = fmt.Sprintf("%s~service:cfg/%s~%s",
      cfg.CFG_SOURCE_ETC,
      dSvcName,
      cfg.CFG_CODER_YAML)
  default:
    uri = fmt.Sprintf("%s~%s~%s",
      cfg.CFG_SOURCE_KUBE,
      dSvcName,
      cfg.CFG_CODER_STRUCT)

  }

  _c, err := cfg.New(uri)
  if err != nil {
    return nil, err
  }
  _svcCfgs.svcs[svc] = _c
  return _c, nil
}

// lookupNetFromCfg ...
func lookupNetFromCfg(svcCfg cfg.Config) (*Net, error) {
  var port string
  err := svcCfg.Scan(CFG_RPC_PORT_PATH, &port)
  if err != nil {
    return nil, err
  }
  var net string
  err = svcCfg.Scan(CFG_RPC_NET_PATH, &net)
  if err != nil {
    return nil, err
  }
  return &Net{Port: port, Net: net}, nil
}

// getDeployedSeviceName ...
func getDeployedSeviceName(svc string) string {
  if svc == "@@self" {
    return cfg.Service
  }

  // get real service name from annotation of cfg
  return svc[2:]
}
