package svr

import (
  "github.com/yoophox/emei/env"
  "github.com/yoophox/emei/flag"
)

func init() {
  if flag.IsHelper() {
    return
  }

  _rootEnv = env.New(nil)
}
