package cfg

import (
  "errors"
  "io/fs"
  "log"
  "os"
  "path"
  "strings"

  "github.com/yolksys/emei/cfg/coder"
  "github.com/yolksys/emei/cfg/source"
  "github.com/yolksys/emei/cmd"
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
  svcName = cmd.String("service", "service name", "")
  if svcName != "" {
    Service = svcName
  }

  err = getCfgFilePath(appPath)
  if err != nil {
    log.Print("cfg init-getCfgFilePath failed: ", err)
    // os.Exit(1)
    return
  }

  s, err := source.Load("local", _path)
  utils.AssertErr(err)

  v, err := coder.Encode("cfg", s)
  utils.AssertErr(err)

  _cfg = &config{
    v: v,
  }

  if svcName == "" {
    GetCfgItem("service", svcName)
    if svcName != "" {
      Service = svcName
    }
  }
}

// 配置文件搜索目录：
// --cfg-file=path
// /etc/bin-name.cfg
// bin-path/cfg/.cfg
func getCfgFilePath(appPath string) error {
  cmdLineFile := ""
  cmdLinePara := "--cfg-file="
  for _, v := range os.Args {
    if strings.Index(v, cmdLinePara) == 0 {
      cmdLineFile = v[len(cmdLinePara):]
      break
    }
  }

  if len(cmdLineFile) > 0 {
    if cmdLineFile[0]^'/' == 0 {
      _path = cmdLineFile
      return nil
    }
    curPath, err := os.Getwd()
    if err != nil {
      return err
    }

    if curPath[len(curPath)-1] != '/' {
      curPath += "/"
    }

    _path = curPath
    return nil
  }

  appName := path.Base(os.Args[0])
  ///~/.config/appName/appName.cfg
  _path = "/~/.config/" + appName + "/" + appName + ".cfg"
  _, err := os.Stat(_path)
  if err == nil {
    return nil
  }

  ///etc
  _path = "/etc/" + appName + ".cfg"
  _, err = os.Stat(_path)
  if err == nil {
    return nil
  }

  // dir of app bin
  fInfo, err := os.Lstat(appPath)
  if err != nil {
    return err
  }

  if fInfo.Mode()&fs.ModeSymlink != 0 {
    appPath, err = os.Readlink(appPath)
    if err != nil {
      return err
    }
  }

  appPath = strings.ReplaceAll(appPath, "\\", "/")
  _path = path.Dir(appPath) + "/cfg/.cfg"
  _, err = os.Stat(_path)
  if err == nil {
    return nil
  }
  //_path = "D:/yolk/code/workspace/go/cfg/.cfg"
  return errors.New("have no cfg file")
}

var (
  _path string
  _cfg  *config
  _uid  string
)

var Service string
