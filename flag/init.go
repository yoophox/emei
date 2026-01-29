package flag

import (
  "os"
  "regexp"
)

func parseCmdLine() {
  _cmdReg := regexp.MustCompile(`^--?([a-zA-Z0-9]+([\.\-_][a-zA-A0-9]+)*)(=(.+))?$`)
  //_sreg := regexp.MustCompile(`^[a-zA-Z]$`)
  for i := 1; i < len(os.Args); i++ {
    strs := _cmdReg.FindAllStringSubmatch(os.Args[i], -1)
    if strs == nil {
      panic("error cmd line parameter:" + os.Args[i])
    }

    if strs[0][4] == "" && i < len(os.Args)-1 && os.Args[i+1][0] != '-' {
      _cmds[strs[0][1]] = os.Args[i+1]
      i++
    } else {
      _cmds[strs[0][1]] = strs[0][4]
    }
  }
}

func init() {
  parseCmdLine()
}
