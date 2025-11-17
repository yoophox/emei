package kube

import (
  "context"
  "fmt"
  "math/bits"
  "math/rand"
  "strings"
  "sync"

  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/cmd"
  "github.com/yolksys/emei/kube/resolver"
)

var (
  Lookup    = lookup
  GetRpcNet = getRpcNet
  GetEPTs   func(svc string) (ips []string, err error) // get all endpoints of a service
)

func lookupLocal(svc, rpc string) (trans string, ip string, port string, err error) {
  ips_, ok := _hosts[svc]
  if !ok {
    err = fmt.Errorf("fail:lookupLocal, mag:have no cfg for svc:%s", svc)
    return
  }
  pos := rand.Intn(len(ips_))
  if pos == _hostsTypLen[svc+"ipv4"] {
    pos++
  }
  ip = ips_[pos]

  ipSec := strings.Split(ip, ";")
  if len(ipSec) > 3 {
    ip = ipSec[0]
    var tp_ []string

    switch rpc {
    case "nrpc":
      tp_ = strings.Split(ipSec[1], ":")
    case "grpc":
      tp_ = strings.Split(ipSec[2], ":")
    case "web":
      tp_ = strings.Split(ipSec[3], ":")
    }
    trans = tp_[0]
    port = tp_[1]
    return
  }

  trans, port, err = getRpcNet(svc, rpc)
  return
}

// lookup ...
// use kube/dns/etc
func lookup(svc, rpc string) (trans string, ip string, port string, err error) {
  trans, port, err = getRpcNet(svc, rpc)
  if err != nil {
    return
  }

  ips, err := resolver.Resover.LookupAddr(context.Background(), svc)
  if err != nil {
    return
  }

  ip = ips[rand.Intn(len(ips))]
  return
}

// getRpcCfgFromeCmdline ...
func getRpcCfgFromCmdline(svc, rpc string) (trans string, port string, err error) {
  switch rpc {
  case "web":
    return _webTrans, _webPort, nil
  }

  return
}

// use kube/etc/local
func getRpcNet(svc, rpc string) (trans string, port string, err error) {
  cfg_, err := getSvcCfg(svc)
  if err != nil {
    return
  }

  err = cfg_.Scan("metadata.labels."+rpc+"-trans", &trans)
  if err != nil {
    return
  }

  allPorts := make([]struct {
    Name string `json:"name,omitempty"`
    Port string `json:"port,omitempty"`
  }, 0)
  err = cfg_.Scan("spec.ports[*]", &allPorts)
  if err != nil {
    return
  }
  for _, value := range allPorts {
    if value.Name == rpc {
      port = value.Port
      break
    }
  }

  return
}

// getSvcCfg ...
func getSvcCfg(svc string) (cfg.Config, error) {
  _svcCfgs.m.RLock()
  cfg_, ok := _svcCfgs.svcs[svc]
  _svcCfgs.m.RUnlock()
  if ok {
    return cfg_, nil
  }

  _svcCfgs.m.Lock()
  defer _svcCfgs.m.Unlock()
  cfg_, ok = _svcCfgs.svcs[svc]
  if ok {
    return cfg_, nil
  }
  cfg_, err := cfg.New(fmt.Sprintf("%s/%s%s~%s", _svcCfgUriPrefix, svc, _svcFileExt, _svcCfgFmt))
  if err != nil {
    return nil, fmt.Errorf("fail:new cfg, par:{%s}, svc:%s, prefix:%s, fmt:%s",
      err.Error(), svc, _svcCfgUriPrefix, _svcCfgFmt)
  }
  _svcCfgs.svcs[svc] = cfg_
  return cfg_, nil
}

// initLocalServers ...
func initLocalServers() {
  localSvr, err := cfg.New(fmt.Sprintf("local~%s~cfg", _localServersPath))
  if err != nil {
    panic(err.Error())
  }

  localSvr.Scan("hosts", &_hosts)
  for key, value := range _hosts {
    i := 0
    for ; i < len(value); i++ {
      if value[i] != `.` {
        continue
      }

      break
    }

    _hostsTypLen[key+"ipv4"] = i // i - i>>bits.TrailingZeros8(uint8(i))&1 // when i > 0, i - 1
    ip6L := len(value) - i
    _hostsTypLen[key+"ipv6"] = ip6L - ip6L>>bits.TrailingZeros8(uint8(ip6L))&1
    _hostsTypLen[key+"all"] = len(value)
  }
}

func init() {
  _grpcTrans = cmd.String("grpctrans", "", "")
  _grpcPort = cmd.String("grpcport", "", "")
  _nrpcTrans = cmd.String("nrpctrans", "", "")
  _nrpcPort = cmd.String("nrpcport", "", "")
  _webTrans = cmd.String("webtrans", "", "")
  _webPort = cmd.String("webport", "", "")

  dns := cmd.String("dns", "", "")       // addr,addr;com or addr,etc or local:path
  svcCfg := cmd.String("svccfg", "", "") // etc or cmd or local:[/]path to your dir

  if svcCfg == "cmd" {
    GetRpcNet = getRpcCfgFromCmdline
  } else if len(svcCfg) > 6 && (svcCfg)[:6] == "local:" {
    // GetRpcCfg = getRpcCfgFromLocal
    _svcCfgUriPrefix = "local~" + strings.TrimSuffix((svcCfg)[6:], "/")
    _svcFileExt = ".yaml"
    // _localSvcCfgPath = *localSvcsCfgPath
  } else if svcCfg == "etc" {
    // GetRpcCfg = getRpcCfgFromEtc
    _svcCfgUriPrefix = "etc~service:cfg/"
  } else {
    //    GetRpcCfg = getRpcCfg
    _svcCfgUriPrefix = "kube~api"
  }

  if len(dns) > 6 && dns[:6] == "local:" {
    // Lookup = lookupDnsServer
    Lookup = lookupLocal
    _localServersPath = dns[6:]
    initLocalServers()
  } else if len(dns) > 0 {
    var domain string
    t := strings.Split(dns, ";")
    if len(t) > 0 {
      domain = t[1]
    }
    resolver.Init(t[0], domain)
  } else {
    // Lookup = lookup
  }
}

type svcCache struct {
  m    sync.RWMutex
  svcs map[string]cfg.Config
}

var (
  _svcCfgs          = svcCache{svcs: map[string]cfg.Config{}}
  _svcCfgUriPrefix  string
  _svcCfgFmt        = "yaml"
  _svcFileExt       = ""
  _grpcTrans        string
  _grpcPort         string
  _nrpcTrans        string
  _nrpcPort         string
  _webTrans         string
  _webPort          string
  _localServersPath string
  _localSvcCfgPath  string
  _hosts            = map[string][]string{} // key = service name, value = ips
  _hostsTypLen      = map[string]int{}      // key = svc name + rr.Type.string(), value = ip num
)
