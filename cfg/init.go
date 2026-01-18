package cfg

import (
  "os"
  "path"

  "github.com/yolksys/emei/cla"
  "github.com/yolksys/emei/utils"
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
