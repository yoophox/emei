package flag

import (
  "flag"
  "os"
)

// NewFlagSet ...
func NewFlagSet(name string) *FlagSet {
  fs_ := &FlagSet{
    FlagSet: flag.NewFlagSet(name, flag.ContinueOnError),
  }

  fs_.args = []string{} // getArgsFromOs(flags...)
  _flagsets[name] = fs_
  return fs_
}

// IsHelper ...
func IsHelper() bool {
  _, oh := _cmds["h"]
  _, ohh := _cmds["help"]
  if !oh && !ohh {
    return false
  }

  return true
}

// Usage ...
func Usage() {
  _, oh := _cmds["h"]
  _, ohh := _cmds["help"]
  if !oh && !ohh {
    return
  }

  // fmt.Println(len(_flagsets), "************")
  // for _, f := range _flagsets {
  //   fmt.Println("**************name", f.Name())
  //   if f.Name() == "logger" {
  //     continue
  //   }
  //   f.PrintDefaults()
  // }

  os.Exit(0)
}
