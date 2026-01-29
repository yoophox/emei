package kube

import (
  "strings"

  "github.com/yoophox/emei/cfg"
  "github.com/yoophox/emei/flag"
)

func init() {
  fs_ := flag.NewFlagSet("kube")
  _etc = fs_.Bool("kube.etc", false, "use etc for lookup")
  _localDir = fs_.String("kube.local", "", "use local file for lookup")
  *_localDir = strings.TrimSuffix(*_localDir, "/")
  err := fs_.Parse()
  if err == flag.ErrHelp {
    return
  }

  uri := ""
  if *_etc {
    LookupServer = lookupServerInEtc
    LookupNet = lookupNetInEtc
    LookupIP = lookupIPInEtc
    LookupEPTs = lookupEPTsInEtc
    uri = cfg.BuildCfgURI(cfg.CFG_SOURCE_ETC, "service:cfg/"+cfg.Service, cfg.CFG_CODER_YAML)
  } else if *_localDir != "" {
    LookupServer = lookupServerInLocal
    LookupNet = lookupNetInLocal
    LookupIP = lookupIPInLocal
    LookupEPTs = lookupEPTsInLocal
    uri = cfg.BuildCfgURI(cfg.CFG_SOURCE_LOCAL, *_localDir+"/"+cfg.Service+".yaml")
  } else {
    LookupServer = lookupServer
    LookupNet = lookupNet
    LookupIP = lookupIP
    LookupEPTs = lookupEPTs
    uri = cfg.BuildCfgURI(cfg.CFG_SOURCE_KUBE, cfg.Service, cfg.CFG_CODER_STRUCT)
  }

  _selfSvcCfg, err = cfg.New(uri)
  if err != nil {
    panic(err)
  }
}
