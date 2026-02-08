package uuid

import (
  "context"
  "fmt"
  "math/rand"
  "os"
  "os/signal"
  "strconv"
  "strings"
  "sync"
  "syscall"
  "time"

  "github.com/sony/sonyflake/v2"
  "github.com/yoophox/emei/cfg"
  "github.com/yoophox/emei/dcs"
  "github.com/yoophox/emei/flag"
  "github.com/yoophox/emei/utils"
)

func init() {
  fs_ := flag.NewFlagSet("uuid")
  _mId := fs_.Uint("uuid.machineid", 0, "machineid for sonyflake uuid generator: 100 < id > 1 << 16")
  err := fs_.Parse()
  if err == flag.ErrHelp {
    return
  }

  if *_mId != 0 && (*_mId <= 100 || *_mId > 1<<16) {
    panic("0 < _machineid < 1 << 16")
  }

  if *_mId == 0 {
    var pnum int
    pnum, err = getPodOrdinal()
    if err == nil {
      *_mId = uint(pnum) + 100
    }
  }

  if *_mId == 0 {
    machineidFromValkey()
    go sigupdate()
  }

  if *_mId > 100 {
    _machineId = uint16(*_mId)
  }

  _localIdGener, err = sonyflake.New(
    sonyflake.Settings{
      MachineID: func() (int, error) {
        return int(_machineId), nil
      },
    },
  )
  if err != nil {
    utils.AssertErr(err)
  }

  go worker()
}

// Get the pod's unique index (e.g., 0, 1, 2)
func getPodOrdinal() (int, error) {
  hostname := os.Getenv("HOSTNAME")
  pos := strings.Index(hostname, "-")
  if pos <= 0 {
    return 0, fmt.Errorf("hostname:"+hostname, "")
  }
  return strconv.Atoi(hostname[pos+1:])
}

// sigupdate ...
func sigupdate() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGHUP)
  for {
    <-c
    machineidFromValkey()
    worker()
  }
}

func machineidFromValkey() {
  _machineidMut.Lock()
  defer _machineidMut.Unlock()

  if _machineId > 100 {
    return
  }

  if !dcs.IsValkeyAct() {
    return
  }

  ctx, cancel := context.WithTimeout(context.Background(), cfg.SysTimeout*time.Second)
  aofSync, errc := dcs.Valkey.ConfigGet(ctx, []string{"appendfsync"})
  cancel()
  if errc != nil || (aofSync["appendfsync"] != "everysec" && aofSync["appendfsync"] != "always") {
    return
  }

  for {
    u, err := doGetFromValkeyargs()
    if err != nil {
      return
    }

    if u == 0 {
      continue
    }

    _machineId = u
  }
}

// doGetFromValkey ...
func doGetFromValkeyargs() (uint16, error) {
  u := rand.Intn(65000)
  if u <= 100 {
    u += 100
  }

  mp_ := dcs.CompriseMachineidPath(uint16(u))
  ctx, cancel := context.WithTimeout(context.Background(), (cfg.SysTimeout+3)*time.Second)
  defer cancel()
  ok, err := dcs.Valkey.MSetNX(ctx, map[string]string{mp_: "a"})
  if err != nil {
    return 0, err
  }
  if !ok {
    return 0, nil
  }

  time.Sleep(3 * time.Second)
  exp, err := dcs.Valkey.ExpireTime(ctx, "")
  if err != nil {
    return 0, err
  }
  if exp != -1 {
    return 0, nil
  }

  go func(mp_ string) {
    for {
      ctx, cancel = context.WithTimeout(context.Background(), cfg.SysTimeout*time.Second)
      ok, err := dcs.Valkey.Expire(ctx, mp_, _vkmachineidExp)
      cancel()
      if err != nil || !ok {
        _machineId = 1
        go machineidFromValkey()
        return
      }

      time.Sleep(_vkMachineidSleep)
    }
  }(mp_)

  return uint16(u), nil
}

var (
  _localIdGener     *sonyflake.Sonyflake
  _machineId        uint16 = 1
  _machineidMut            = &sync.Mutex{}
  _uuidCh                  = make(chan int64, 100)
  _uuidWorking      bool   = false
  _vkmachineidExp          = 5*time.Hour + 10*time.Second
  _vkMachineidSleep        = 5 * time.Hour
)
