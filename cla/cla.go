package cla

import (
  "fmt"
  "os"
  "regexp"
  "strings"

  "github.com/yoophox/emei/utils"
)

func Bool(name, desc string, dft bool, short ...rune) bool {
  var s string = ""
  if len(short) == 1 {
    s = string(short[0])
  }

  lv_, lok := _cmdDesc[name]
  sv_, sok := _cmdDesc[s]
  if lok && sok {
    _reduplicativeCmd = append(_reduplicativeCmd, fmt.Sprintf("%s, %s", name, s))
    return false
  }

  var v *cdesc

  if lok {
    v = lv_
  } else if sok {
    v = sv_
  } else {
    return dft
  }

  // v.typ = "bool"
  v.used = true

  switch v.value {
  case "", "true":
    return true
  case "false":
    return false
  default:
    _errCmds = append(_errCmds, fmt.Sprintf("bool:%s:%s", name, v.value))
  }

  return false
}

func String(name, desc string, dft string, short ...rune) string {
  var s string = ""
  if len(short) == 1 {
    s = string(short[0])
  }

  lv_, lok := _cmdDesc[name]
  sv_, sok := _cmdDesc[s]
  if lok && sok {
    _reduplicativeCmd = append(_reduplicativeCmd, fmt.Sprintf("%s, %s", name, s))
    return ""
  }

  var v *cdesc

  if lok {
    v = lv_
  } else if sok {
    v = sv_
  } else {
    return dft
  }

  v.used = true
  return v.value
}

func Int64(name, desc string, dft int64, short ...rune) int64 {
  return dft
}

func Uint64(name, desc string, dft uint64, short ...rune) uint64 {
  return dft
}

// Lookup ...
func Lookup(key string) (string, bool) {
  return "", false
}

// Assert ...
func Assert() {
  _, isHelp := _cmdDesc["help"]
  _, isH := _cmdDesc["h"]
  if isH || isHelp {
    usage()
    os.Exit(0)
  }

  noAllUsed := false
  for _, value := range _cmdDesc {
    if !value.used {
      noAllUsed = true
      _notSupported = append(_notSupported, value.name)
      // break
    }
  }

  if len(_errCmds) == 0 && len(_reduplicativeCmd) == 0 && !noAllUsed {
    return
  }

  usage()
  utils.AssertErr(fmt.Errorf("cmd error"))
}

// parseCmd ...
func parseCmd() {
  _reg := regexp.MustCompile(`^[a-zA-Z\.]+(=.+)?$`)
  for i := 1; i < len(os.Args); i++ {
    if len(os.Args[i]) > 2 && os.Args[i][:2] == "--" && _reg.MatchString(os.Args[i][2:]) {
      pos := strings.Index(os.Args[i][2:], "=")
      if pos > 0 {
        _cmdDesc[os.Args[i][2:2+pos]] = &cdesc{name: os.Args[i][2:pos], value: os.Args[i][2+pos+1:]}
      } else if i < len(os.Args)-1 && os.Args[i+1][0] != '-' {
        _cmdDesc[os.Args[i][2:]] = &cdesc{name: os.Args[i][2:], value: os.Args[i+1]}
        i++
      } else {
        _cmdDesc[os.Args[i][2:]] = &cdesc{name: os.Args[i][2:]}
      }
    } else if len(os.Args[i]) == 2 && os.Args[i][0] == '-' && _reg.MatchString(os.Args[i][2:]) {
      if i < len(os.Args)-1 && os.Args[i+1][0] != '-' {
        _cmdDesc[os.Args[i][1:]] = &cdesc{value: os.Args[i+1]}
        i++
      } else {
        _cmdDesc[os.Args[i][1:]] = &cdesc{}
      }
    } else {
      _errCmds = append(_errCmds, os.Args[i])
    }
  }
}

// check ...
// func check(name, desc, typ, d string) bool {
//   f, fc_, l := utils.GetCallInfo(2)
//   cio := fmt.Sprintf("%s-%s-%d", f, fc_, l)
//   descv, ok := _cmdDesc[name]
//   if ok {
//     _errCmds = append(_errCmds, fmt.Sprintf("fail:%s, reason:reregister, origin:%s, cur:%s", typ, descv.reger, cio))
//     return false
//   }
//
//   _cmdDesc[name] = &cdesc{
//     desc:  desc,
//     short: d,
//     typ:   typ,
//     reger: cio,
//   }
//   return true
// }

// usage ...
func usage() {
  if len(_reduplicativeCmd) > 0 {
    fmt.Print("redupicative parameters:\n")
    for _, value := range _reduplicativeCmd {
      fmt.Print("    ", value)
    }
  }

  if len(_errCmds) > 0 {
    for _, value := range _errCmds {
      fmt.Print("error cmds:\n")
      fmt.Print("    ", value)
    }
  }

  if len(_notSupported) > 0 {
    fmt.Print("not supportted cmds:\n")
    for _, value := range _notSupported {
      fmt.Print("    ", value)
    }
  }

  fmt.Print("usage:\n")
  for k, v := range _cmdDesc {
    fmt.Print("    ", k, "    ", v)
  }
}

func init() {
  parseCmd()
}

// argument value
type cdesc struct {
  name  string
  short string
  desc  string
  value string
  typ   string
  used  bool
  // reger string
  v any
}

var (
  _cmdDesc          = map[string]*cdesc{}
  _reduplicativeCmd = []string{} //"long, short"
  _errCmds          = []string{}
  _notSupported     = []string{}
)
