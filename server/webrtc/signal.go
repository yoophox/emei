package webrtc

import (
  "encoding/json"
  "fmt"
  "io"
  "reflect"

  "github.com/gorilla/websocket"
  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/env"
)

type Webrtc0 struct{}

// called signal server by mdn
func (w *Webrtc0) Signal(e env.Env, conn any) {
  defer e.Return()

  newCli := func(uid string) {
    if uid == "" {
      e.AssertErr(fmt.Errorf("fail:signal login, reason:null uid"))
    }

    cli := &client{
      uid: uid,
      oto: []*otochan{},
      otm: map[string]*otmroom{},
      mtm: map[string]*mtmroom{},
    }

    _, ok := _clients[uid]
    if ok {
      e.AssertErr(fmt.Errorf("fail:signal login, reason:repeated uid, uid:%s", uid))
    }

    _clients[uid] = cli
  }

  alloc := func(m *msg) {
    switch m.Class {
    case SigClassOTOCall:
    case SigClassOTOCallee:
    case SigClassAllocOTM:
    case SigClassAllocMTM:
    case SigClassAllocLCT:
    }
  }

  // pass to group
  passg := func(m *msg) error {
    var peerid string
    var ok bool

    switch m.Class {
    case SigClassOTOCallee, SigClassAllocCli:
      peerid, ok = m.Attrs["target"]
      if !ok {
        return fmt.Errorf("fail:signal passg, reason:no target, class:%d", m.Class)
      }
    case SigClassCfg:
      peerid, ok = m.Attrs["name"]
      if !ok {
        return fmt.Errorf("fail:signal passg, reason:no peer name, class:%d", m.Class)
      }
    }

    group, ok := _groups[peerid]
    if !ok {
      return fmt.Errorf("fail:signal passg, reason:not alloc for peer, peername:%s, class:%d", peerid, m.Class)
    }

    return group.hSig(m)
  }

  var response func(v any) error
  var decode func(v any) error
  isws := false
  var m msg
  switch c := conn.(type) {
  case io.ReadWriter:
    dec := json.NewDecoder(c)
    decode = dec.Decode
    enc := json.NewEncoder(c)
    response = enc.Encode

  case *websocket.Conn:
    isws = true
    decode = c.ReadJSON
    response = c.WriteJSON
  default:
    e.AssertErr(fmt.Errorf("fail:signal, reason: err conn type, conn type:%+v", reflect.TypeOf(conn).Name()))
  }

  var retv msg
  retv.Class = SigClassResponse
  for {
    err := decode(&m)
    e.AssertErr(err)
    if m.Class == SigClassLogin {
      var uid string
      if isws {
        uid = cfg.GetUID()
      } else {
        uid, _ = m.Attrs["uid"]
      }
      newCli(uid)
    } else if m.Class == SigClassAllocLCT ||
      m.Class == SigClassAllocMTM ||
      m.Class == SigClassAllocOTM ||
      m.Class == SigClassOTOCall {
      alloc(&m)
    } else if m.Class == SigClassCfg ||
      m.Class == SigClassOTOCallee ||
      m.Class == SigClassAllocCli {
      err = passg(&m)
    } else {
      err = fmt.Errorf("fail:signal, reason:not support class, class:%d", m.Class)
    }

    if err != nil {
      defer func() {
        delete(retv.Attrs, "code")
        delete(retv.Attrs, "reason")
      }()
      retv.Attrs["code"] = "1"
      retv.Attrs["reason"] = err.Error()
    }
    response(&retv)
  }
}

type msg struct {
  Class byte              `json:"class,omitempty"`
  Typ   byte              `json:"type,omitempty"`
  Attrs map[string]string `json:"attrs,omitempty"`
}

type client struct {
  uid string
  oto []*otochan // two peer all are webrtc then pass data by turn
  otm map[string]*otmroom
  mtm map[string]*mtmroom
  lct map[string]*lctroom
}

var (
  _peers   = map[string]*peer{} // key==peername
  _groups  = map[string]group{} // key==peername, value= *otochan, *otmroom, *mtmroom, *lctroom
  _clients = map[string]*client{}
  // _oto   = map[string]*otochan{}
  // _otm   = map[string]*otmroom{}
  // _mtm   = map[string]*mtmroom{}
  // _lct   = map[string]*lctroom{}
)
