package dns

import (
  "context"
  "math/rand"
  "sync"

  "github.com/yolksys/emei/etc/etcintra"
  "github.com/yolksys/emei/etc/etcout"
  "github.com/yolksys/emei/utils"
)

// Start ...
func Start() {
  if _exportIp != "" {
    etcintra.EtcCli.Put(context.Background(), "service:ip", _service, utils.HostId(), _exportIp,
      etcintra.WithLeaseIdOption(_lease.GetID()))
    etcintra.EtcCli.KeepLease(_lease.GetID())
  }
}

// GetSvcIp ...
func GetSvcIp(svc string) (string, error) {
  _serviceIpsMux.Lock()
  defer _serviceIpsMux.Unlock()
  ips, ok := _serviceIps[svc]
  if ok {
    return ips[rand.Intn(len(ips))], nil
  }

  ips_, err := etcintra.EtcCli.Get("service:ip", svc)
  if err != nil {
    return "", err
  }
  for _, value := range ips_ {
    ips = make([]string, 0, len(ips_))
    ips = append(ips, value)
  }
  _serviceIps[svc] = ips
  _, err = etcintra.EtcCli.Watch("service:ip", svc)
  if err != nil {
    return "", err
  }

  return ips[rand.Intn(len(ips))], nil
}

func Init(export, svc string) {
  _exportIp = export
  _service = svc
}

// recvIpEvent ...
func recvIpEvent() {
  for {
  }
}

var (
  _lease         etcintra.Lease
  _ipWatchChan   = make([]<-chan etcout.Event, 0, 10)
  _serviceIps    = make(map[string][]string)
  _serviceIpsMux = sync.Mutex{}
  _exportIp      string
  _service       string
)
