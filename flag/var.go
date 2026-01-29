package flag

import "flag"

var ErrHelp = flag.ErrHelp

var (
  _flagsets = map[string]*FlagSet{}
  _cmds     = map[string]string{}
)
