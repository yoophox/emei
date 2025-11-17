package cmd

import (
  "fmt"
  "os"
  "regexp"
  "strings"

  "github.com/yolksys/emei/utils"
)

func Bool(name, desc string, dft bool, short ...rune) bool {
  var s string = ""
  if len(short) == 1 {
    s = string(short[0])
  }

  ok := check(name, desc, "bool", s)
  if !ok {
    return false
  }

  lv_, lok := _cmds[name]
  sv_, sok := _cmds[s]
  if lok && sok {
    _reduplicativeCmd = append(_reduplicativeCmd, fmt.Sprintf("%s, %s", name, s))
    return false
  }

  var v *value

  if lok {
    v = lv_
  } else if sok {
    v = sv_
  } else {
    return dft
  }
  v.used = true

  if v.str == "" || v.str == "true" {
    return true
  } else if v.str == "false" {
    return false
  } else {
    _errCmds = append(_errCmds, fmt.Sprintf("bool:%s:%s", name, v.str))
  }

  return false
}

func String(name, desc string, dft string, short ...rune) string {
  var s string = ""
  if len(short) == 1 {
    s = string(short[0])
  }

  ok := check(name, desc, "bool", s)
  if !ok {
    return ""
  }

  lv_, lok := _cmds[name]
  sv_, sok := _cmds[s]
  if lok && sok {
    _reduplicativeCmd = append(_reduplicativeCmd, fmt.Sprintf("%s, %s", name, s))
    return ""
  }

  var v *value

  if lok {
    v = lv_
  } else if sok {
    v = sv_
  } else {
    return dft
  }
  v.used = true

  return v.str
}

func Int64(name, desc string, short rune, dft int64) int64 {
  return dft
}

func Uint64(name, desc string, short rune, dft uint64) uint64 {
  return dft
}

// Lookup ...
func Lookup(key string) (string, bool) {
  return "", false
}

// Assert ...
func Assert() {
  noAllUsed := false
  for _, value := range _cmds {
    if !value.used {
      noAllUsed = true
      break
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
  reg := regexp.MustCompile(`^[a-zA-Z\.]+(=.+)?$`)
  for i := 1; i < len(os.Args); i++ {
    if len(os.Args[i]) > 2 && os.Args[i][:2] == "--" && reg.MatchString(os.Args[i][2:]) {
      pos := strings.Index(os.Args[i][2:], "=")
      if pos > 0 {
        _cmds[os.Args[i][2:2+pos]] = &value{os.Args[i][2+pos+1:], false}
      } else if i < len(os.Args)-1 && os.Args[i+1][0] != '-' {
        _cmds[os.Args[i][2:]] = &value{os.Args[i+1], false}
        i++
      } else {
        _cmds[os.Args[i][2:]] = &value{}
      }
    } else if len(os.Args[i]) == 2 && os.Args[i][0] == '-' && reg.MatchString(os.Args[i][2:]) {
      if i < len(os.Args)-1 && os.Args[i+1][0] != '-' {
        _cmds[os.Args[i][1:]] = &value{os.Args[i+1], false}
        i++
      } else {
        _cmds[os.Args[i][1:]] = &value{}
      }
    } else {
      _errCmds = append(_errCmds, os.Args[i])
    }
  }
}

// check ...
func check(name, desc, typ, d string) bool {
  f, fc_, l := utils.GetCallInfo(2)
  cio := fmt.Sprintf("%s-%s-%d", f, fc_, l)
  descv, ok := _cmdDesc[name]
  if ok {
    _errCmds = append(_errCmds, fmt.Sprintf("fail:%s, reason:reregister, origin:%s, cur:%s", typ, descv.reger, cio))
    return false
  }

  _cmdDesc[name] = &cdesc{
    desc:  desc,
    short: d,
    typ:   typ,
    reger: cio,
  }
  return true
}

// usage ...
func usage() {
}

func init() {
  parseCmd()
}

type value struct {
  str  string
  used bool
}

type cdesc struct {
  desc  string
  short string
  typ   string
  reger string
}

var (
  _cmds             = map[string]*value{}
  _cmdDesc          = map[string]*cdesc{}
  _reduplicativeCmd = []string{} //"long, short"
  _errCmds          = []string{}
)
