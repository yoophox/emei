package dcs

import (
  "context"
  "fmt"
  "os"
  "os/signal"
  "strconv"
  "syscall"
  "time"

  glide "github.com/valkey-io/valkey-glide/go/v2"
  "github.com/valkey-io/valkey-glide/go/v2/config"
  "github.com/yoophox/emei/cfg"
  "github.com/yoophox/emei/kube"
  "github.com/yoophox/emei/names"
)

// IsValkeyAct ...
func IsValkeyAct() bool {
  return _valkeyAct
}

func init() {
  configValkey()
  go sigupdate()
}

// update ...
func sigupdate() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGHUP)
  for {
    <-c
    configValkey()
  }
}

// reconfigValkey ...
func configValkey() {
  if Valkey != nil {
    // Valkey.Close()
    // Valkey = nil
    return
  }
  svr, err := kube.LookupServer(names.NAME_SERVICE_VALKEY)
  if err != nil {
    fmt.Println("********   info: lookup valkey server error, ", err.Error(), "   *********")
    return
  }

  ccf := config.NewClusterClientConfiguration()
  node := &config.NodeAddress{}
  node.Host = svr.IP
  port, _ := strconv.ParseInt(svr.Net.Ports[names.NAME_SERVICE_PORT_VALKEY_CLIENT].Port, 10, 32)
  node.Port = int(port)
  ccf.WithAddress(node)
  Valkey, err = glide.NewClusterClient(ccf)
  if err != nil {
    fmt.Println("********   info: glide newclusterclient error, err:", err.Error(), "   *********")
    return
  }
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()
  _, err = Valkey.Ping(ctx)
  if err != nil {
    fmt.Println("********   info: glide ping error, err:", err.Error(), "   *********")
    return
  }

  _valkeyAct = true
  registerSvcInroot()
}

// registerSvcInroot ...
func registerSvcInroot() {
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  Valkey.SAdd(ctx, _ROOT, []string{cfg.Service})
}

var (
  _valkeyAct                      = false
  Valkey     *glide.ClusterClient = nil
)
