package dcs

import (
  "fmt"

  "github.com/yoophox/emei/cfg"
)

// CompriseMachineidPath ...
func CompriseMachineidPath(id uint16) string {
  return fmt.Sprintf("%s.machineid.%d", cfg.Service, id)
}

// _root = {services}
// path = svcname.class.
var _ROOT = "/"
