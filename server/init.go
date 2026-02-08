package svr

import (
  "github.com/yoophox/emei/env"
  "github.com/yoophox/emei/flag"
)

func init() {
  fs_ := flag.NewFlagSet("server")
  ori := fs_.String("web.ori", "", "cors origin default is all")
  _port = fs_.String("port", "443", "rpc and web port")
  _host = fs_.String("host", "", "host which listen at, default is all")
  err := fs_.Parse()
  if err == flag.ErrHelp {
    return
  }

  _webCorsOri = []string{*ori}
  _rootEnv = env.New(nil)
}
