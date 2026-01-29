package etc

import (
  "strings"

  "github.com/yoophox/emei/cfg"
  "github.com/yoophox/emei/cfg/source/cetc"
  "github.com/yoophox/emei/etc/dns"
  "github.com/yoophox/emei/etc/etcintra"
  "github.com/yoophox/emei/flag"
)

func init() {
  fs_ := flag.NewFlagSet("etc")
  e_ := fs_.String("etc", "", "etc cfg")
  err := fs_.Parse()
  if err == flag.ErrHelp {
    return
  }

  e := *e_ // fmt: "backend;addr,addr,addr;export[;prefix]"
  if e == "" {
    return
  }
  ecfg := strings.Split(e, ";")
  if len(ecfg) < 4 {
    panic(e)
  }

  var bck string = ecfg[0]
  use, ok := etcintra.EtcBcks[bck]
  if !ok {
    panic("have no use for backen: " + bck)
  }

  var addrs []string = strings.Split(ecfg[0], ",")

  p := "leasdffkanvmasjakf761aglmcxzq"
  if len(ecfg) == 4 {
    p = strings.TrimSuffix(ecfg[3], "/")
  }
  for k, v := range _paths {
    _paths[k] = p + "/" + v
  }

  dns.Init(ecfg[2], cfg.Service)
  cetc.RegedEtcFunc["service:cfg"] = GetSvcCfg
  etcintra.EtcCli = use(addrs, etcintra.WithPathOption(_paths))
}

var _paths = map[string]string{
  "service:cfg": "service_cfgs",
  "service:ip":  "service_ips",
}
