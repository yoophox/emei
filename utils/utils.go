package utils

import (
  "fmt"
  "log"
  "net"
  "os"
  "os/exec"
  "path"
  "runtime"
  "strings"
  "sync"
)

func AssertErr(e error) {
  if e != nil {
    f, fu_, l := GetCallInfo(1)
    ci_ := fmt.Sprintf(", file:%s, func:%s, line:%d", f, fu_, l)
    log.Fatal(e, ci_)
  }
}

// ...
func AssertTrue(a bool) {
  if a {
    f, fu_, l := GetCallInfo(1)
    ci_ := fmt.Sprintf(", file:%s, func:%s, line:%d", f, fu_, l)
    log.Fatal(ci_)
  }
}

// ...
func catch() {
}

// skip = 0, return name of caller of GetCallInfo
func GetCallInfo(skip int) (file string,
  function string, line int,
) {
  pc_, file, line, ok := runtime.Caller(skip + 1)
  if !ok {
    return
  }
  file = path.Base(file)

  function = runtime.FuncForPC(pc_).Name()
  return
}

// GetPanicFrame ...
// skip == 0, return fram of caller 0f GetPanicFrame
func GetPanicFrame(skip int) *runtime.Frame {
  _, f, _ := GetCallInfo(skip + 1)
  if f == "" {
    return nil
  }

  pc_ := ptrPool.Get().([]uintptr)
  n := runtime.Callers(skip+1, pc_)
  frames := runtime.CallersFrames(pc_[:n])
  frame, _ := frames.Next()
  ok := true
  for ; ok; frame, ok = frames.Next() {
    if frame.Function == f {
      break
    }
  }

  return &frame
}

// IpType ...
// return: "ip4" or "ip6"
func IpType(ip string) string {
  if len(strings.Split(ip, ":")) > 1 {
    return "ip6"
  }

  return "ip4"
}

// ...
func HostId() string {
  switch runtime.GOOS {
  case "linux":
    if b, err := os.ReadFile("/etc/machine-id"); err == nil {
      return strings.TrimSpace(string(b))
    }
    if b, err := os.ReadFile("/etc/machine-id"); err == nil {
      return strings.TrimSpace(string(b))
    }
  case "freebsd":
    if b, err := os.ReadFile("/etc/hostid"); err == nil {
      return strings.TrimSpace(string(b))
    }

    if s, err := ExecCommand("kenv", "-q", "smbios.system.uuid"); err == nil {
      return strings.TrimSpace(string(s))
    }
  // case "windows":
  //   k, err := registry.OpenKey(
  //     registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`,
  //     registry.QUERY_VALUE|registry.WOW64_64KEY,
  //   )
  //   if err != nil {
  //     return "", err
  //   }
  //   defer k.Close()

  //   guid, _, err := k.GetStringValue("MachineGuid")
  //   if err != nil {
  //     return ""
  //   }

  //   return guid

  case "darwin":
    result, err := ExecCommand("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
    if err != nil {
      return ""
    }

    lines := strings.Split(result, "\n")
    for _, line := range lines {
      if strings.Contains(line, "IOPlatformUUID") {
        parts := strings.Split(line, " = ")
        if len(parts) == 2 {
          return strings.Trim(parts[1], "\"")
        }
        break
      }
    }

    return ""
  }

  return ""
}

func ExecCommand(name string, arg ...string) (string, error) {
  cmd := exec.Command(name, arg...)
  b, err := cmd.Output()
  if err != nil {
    return "", err
  }

  return string(b), nil
}

// IsIpv4 ...
func IsIpv4(s string) bool {
  i := net.ParseIP(s)
  return i != nil && i.To4() != nil
}

// IsIp ...
func IsIp(s string) bool {
  return net.ParseIP(s) != nil
}

var (
  ptrPool = sync.Pool{
    New: func() any {
      var ptr [64]uintptr
      return ptr[0:64]
    },
  }
  TimeDefaultFmt = "2006-01-02 15:04:05.999999999 -0700 MST"
)
