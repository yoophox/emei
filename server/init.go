package svr

import (
  "github.com/yoophox/emei/env"
  "github.com/yoophox/emei/flag"
)

func init() {
  fs_ := flag.NewFlagSet("server")
  ori := fs_.String("web.ori", "", "cors origin default is all")
  err := fs_.Parse()
  if err == flag.ErrHelp {
    return
  }

  _webCorsOri = []string{*ori}
  _rootEnv = env.New(nil)
}
