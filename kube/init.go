package kube

import (
  "strings"

  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/cla"
)

func init() {
  _etc = cla.String("kube-etc", "etc cfg", "")
  _localDir = cla.String("kube-local", "local cfg", "")
  _localDir = strings.TrimSuffix(_localDir, "/")

  uri := ""
  if _etc != "" {
    LookupServer = lookupServerInEtc
    LookupNet = lookupNetInEtc
    LookupIP = lookupIPInEtc
    LookupEPTs = lookupEPTsInEtc
    uri = cfg.BuildCfgURI(cfg.CFG_SOURCE_ETC, "service:cfg/"+cfg.Service, cfg.CFG_CODER_YAML)
  } else if _localDir != "" {
    LookupServer = lookupServerInLocal
    LookupNet = lookupNetInLocal
    LookupIP = lookupIPInLocal
    LookupEPTs = lookupEPTsInLocal
    uri = cfg.BuildCfgURI(cfg.CFG_SOURCE_LOCAL, _localDir+"/"+cfg.Service+".yaml")
  } else {
    LookupServer = lookupServer
    LookupNet = lookupNet
    LookupIP = lookupIP
    LookupEPTs = lookupEPTs
    uri = cfg.BuildCfgURI(cfg.CFG_SOURCE_KUBE, cfg.Service, cfg.CFG_CODER_STRUCT)
  }

  var err error

  _selfSvcCfg, err = cfg.New(uri)
  if err != nil {
    panic(err)
  }
}
