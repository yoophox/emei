package kube

import (
  "fmt"

  "github.com/yoophox/emei/cfg"
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/kube/errors"
  "github.com/yoophox/emei/names"
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

  if svc == names.NAME_SERVICE_SELF {
    return _selfSvcCfg, nil
  }

  dSvcName := getDeployedSeviceName(svc)
  if dSvcName == "" {
    return nil, errs.ErrorF(errs.ErrId(errors.ERR_ID_KUBE_EMPTY_SERVICE_NAME))
  }
  uri := ""
  switch uriType {
  case "local":
    uri = fmt.Sprintf("%s~%s", cfg.CFG_SOURCE_LOCAL,
      *_localDir+"/"+dSvcName+".yaml")
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
func lookupNetFromCfg(svc string, svcCfg cfg.Config) (*Net, error) {
  net := &Net{svc: svc}

  err := svcCfg.Scan(CFG_SERVICE_PORTS_PATH, &net.ports)
  if err != nil {
    return nil, err
  }
  net.Ports = make(map[string]*ServicePort, len(net.ports))
  for _, p := range net.ports {
    if p.TargetPort == "" {
      p.TargetPort = p.Port
    }
    net.Ports[p.Name] = p
  }

  return net, nil
}

// getDeployedSeviceName ...
func getDeployedSeviceName(svc string) string {
  if svc == names.NAME_SERVICE_SELF {
    return cfg.Service
  }

  var name string
  err := _selfSvcCfg.Scan(CFG_ANNOTATIONS_PRE+svc, &name)
  if err == nil || name != "" {
    return name
  }

  if svc[len(svc)-1] == '@' {
    return ""
  }

  return svc[2:]
}
