package svr

import (
  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/log"
)

// dispatch ...
func dispatchRpc(l *linkTx) {
  defer func() {
    log.Event("link end", "")
  }()

  for {
    var tja env.Tjatse
    err := l.cc.Decode(&tja)
    if err != nil {
      log.Event()
      l.ReadWriteCloser.Close()
      return
    }
    var topic string
    err = l.cc.Decode(&topic)
    if err != nil {
      log.Event()
      l.ReadWriteCloser.Close()
      return
    }

    e := env.New(&tja)
    assist(e, l)
    e.Release()
  }
}

// ...
func assist(e env.Env, l *linkTx) {
  defer func() {
  }()
  defer e.Return()
}
