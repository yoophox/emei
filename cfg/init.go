package cfg

import (
  "os"
  "path"

  "github.com/yoophox/emei/cla"
  "github.com/yoophox/emei/utils"
)

// private
func init() {
  appPath, err := os.Executable()
  if err != nil {
    utils.AssertErr(err)
  }

  Service = path.Base(appPath)
  var svcName string
  svcName = cla.String("service", "service name", "")
  if svcName != "" {
    Service = svcName
  }
}

var Service string
