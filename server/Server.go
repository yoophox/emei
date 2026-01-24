package svr

import (
  "os"
  "os/signal"

  "github.com/yoophox/emei/kube"
  "github.com/yoophox/emei/utils"
)

// Serve ...
func ServeFor(fo SvrFor, rcvr ...any) {
  switch fo {
  case SERVER_FOR_RPC:
    for _, v := range rcvr {
      err := parseRpc(v)
      if err != nil {
        panic(err)
      }
    }
  case SERVER_FOR_WEB:
    for _, v := range rcvr {
      err := parseWeb(v)
      if err != nil {
        panic(err)
      }
    }
  default:
    panic("not support svrfor:" + fo)
  }
}

// Serve ...
func Serve() {
  netx, err := kube.LookupNet("@@self")
  utils.AssertErr(err)
  err = listenQuic(":" + netx.Port)
  utils.AssertErr(err)
  err = listenTcp(":" + netx.Port)
  utils.AssertErr(err)
  wait()
}

// init(arg ...
func init() {
  // for {reger.RegRpc()}
  // for _, value := range rpc.RpcImps {
  //   reger.RegRpc(value)
  // }

  signal.Notify(_sigCh, os.Interrupt)
  go sigHandler()
}

// wait ...
func wait() {
  //_ = log.New(context.Background())
  _rootEnv.Wait()
}

// sigHandler ...
func sigHandler() {
  <-_sigCh
  //_rootEnv
}

var _sigCh = make(chan os.Signal, 1)
