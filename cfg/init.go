package cfg

import (
  "os"
  "path"

  "github.com/yoophox/emei/flag"
  "github.com/yoophox/emei/utils"
)

// private
func init() {
  fs_ := flag.NewFlagSet("cfg")
  sn_ := fs_.String("service", "", "service name")
  err := fs_.Parse()
  if err == flag.ErrHelp {
    return
  }

  appPath, err := os.Executable()
  if err != nil {
    utils.AssertErr(err)
  }

  Service = path.Base(appPath)
  var svcName string

  svcName = *sn_
  if svcName != "" {
    Service = svcName
  }
}

var Service string
