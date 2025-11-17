package etc

import (
  "github.com/yolksys/emei/etc/dns"
  "github.com/yolksys/emei/etc/etcintra"
  "github.com/yolksys/emei/etc/etcout"
)

// Get ...
func GetSvcCfg(svc string) (string, error) {
  t, err := etcintra.EtcCli.Get("service:cfg", "", svc)
  v, _ := t[svc]
  return v, err
}

// Watch ...
func WatchSvcCfg(svc string) (<-chan etcout.Event, error) {
  return etcintra.EtcCli.Watch("service:cfg", "", svc)
}

// GetSvcIp ...
func GetSvcIp(svc string) (string, error) {
  return dns.GetSvcIp(svc)
}

// GetEPTs ...
func GetEPTs(svc string) (ips []string, err error) {
  return nil, nil
}
