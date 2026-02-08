package uuid

import (
  "context"
  "crypto/tls"
  "io"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "time"

  "github.com/quic-go/quic-go"
  "github.com/quic-go/quic-go/http3"
  "github.com/yoophox/emei/kube"
  "github.com/yoophox/emei/names"
  "github.com/yoophox/emei/utils"
)

// worker ...
func worker() {
  if _uuidWorking {
    return
  }
  net, err := kube.LookupNet(names.NAME_SERVICE_UUID)
  if err != nil {
    return
  }

  for range 5 {
    go getter(net)
  }
}

// getter ...
func getter(net *kube.Net) {
  for {
    ip, err := kube.LookupIP(names.NAME_SERVICE_UUID)
    if err != nil {
      time.Sleep(10 * time.Second)
      continue
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    v := &url.Values{}
    v.Add("num", "100")
    url := utils.CompriseAddr(ip, net.GetPortByName(names.NAME_SERVICE_PORT_DEFAULT_QUIC)) + "/uuid.uuids"
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(v.Encode()))
    if err != nil {
      panic(err)
    }
    for {
      // resp body: 16 byte hex uuid seperated by comma
      resp, err := client.Do(req)
      cancel()
      if err != nil {
        break
      }

      hexUUID := make([]byte, 16)
      n, err := io.ReadFull(resp.Body, hexUUID)
      if n < 16 || err != nil {
        resp.Body.Close()
        break
      }

      uuid, err := strconv.ParseInt(string(hexUUID), 16, 64)
      if err != nil {
        resp.Body.Close()
        break
      }

      _uuidWorking = true
      _uuidCh <- uuid
    }
  }
}

var tr = &http3.Transport{
  // set a TLS client config, if desired
  TLSClientConfig: &tls.Config{
    NextProtos:         []string{http3.NextProtoH3}, // set the ALPN for HTTP/3
    InsecureSkipVerify: true,
  },
  QUICConfig: &quic.Config{}, // QUIC connection options
}

// defer tr.Close()
var client = &http.Client{
  Transport: tr,
}
