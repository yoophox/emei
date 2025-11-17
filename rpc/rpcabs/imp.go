package rpcabs

import (
  "crypto/tls"
  "errors"
  "fmt"
  "go/token"
  "io"
  "reflect"
  "strings"
  "sync"

  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/kube"
  "github.com/yolksys/emei/log"
)

func (s *RpcImp) RegRcvr(rcvrs map[string]*Recver) {
  s.Recvs = rcvrs
  return
}

func (s *RpcImp) ErrChan() chan error {
  return s.Err
}

func (s *RpcImp) IsClosed() bool {
  return s.IsClosed_
}

func (s *RpcImp) Init() error {
  // ip, port, Trans, cert
  // init
  s.Recvs = make(map[string]*Recver)
  s.Wg_ = &sync.WaitGroup{}
  s.Err = make(chan error)
  s.IsClosed_ = false
  var err error
  s.Trans, s.Port, err = kube.GetRpcNet(cfg.Service, s.Name)
  if err != nil {
    return err
  }
  err = cfg.GetCfgItem("net.crypto.key", &s.Key)
  if err != nil {
    return err
  }
  err = cfg.GetCfgItem("net.crypto.cert", &s.Cert)
  if err != nil {
    return err
  }

  s.Certificate, err = tls.LoadX509KeyPair(s.Cert, s.Key)
  if err != nil {
    return fmt.Errorf("fail:tls.loadX509keypair, reason:%s", err.Error())
  }

  return nil
}

// parseRcvr ...
func ParseRcvr(rcvr any) (*Recver, error) {
  typ := reflect.TypeOf(rcvr)
  rcvrV := reflect.ValueOf(rcvr)

  sname := reflect.Indirect(rcvrV).Type().Name()
  if sname == "" {
    s := "rpc.Register: no service name for type " + typ.String()
    return nil, errors.New(s)
  }
  if !token.IsExported(sname) {
    s := "rpc.Register: type " + sname + " is not exported"
    return nil, errors.New(s)
  }

  // Install the methods
  mets, eMsg := suitableMethods(typ)
  if len(mets) == 0 {
    str := ""
    // To help the user, see if a pointer receiver would work.
    mets, _ = suitableMethods(reflect.PointerTo(typ))
    if len(mets) != 0 {
      str = "rpc.Register: type " + sname + " (hint: pass a pointer to value of that type)"
    } else {
      str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
    }

    return nil, errors.New(str)
  }

  if len(eMsg) != 0 {
    log.Info("have error met", eMsg)
  }

  _rcvr := &Recver{}
  _rcvr.Val = rcvrV
  _rcvr.Typ = typ
  _rcvr.Mets = mets
  _rcvr.Name = sname

  sname = strings.ToLower(sname)

  return _rcvr, nil
}

// suitableMethods returns suitable Rpc methods of typ. It will log
// errors if logErr is true.
func suitableMethods(typ reflect.Type) (map[string]*MethodType, []string) {
  mets := make(map[string]*MethodType)

  msgErr := make([]string, 0)

  for m := 0; m < typ.NumMethod(); m++ {
    method := typ.Method(m)
    mtype := method.Type
    mname := method.Name
    // Method must be exported.
    if !method.IsExported() {
      continue
    }

    inNums := mtype.NumIn()
    // Method needs two ins: receiver, emei.Env
    if inNums < 2 {
      msgErr = append(msgErr,
        fmt.Sprintf("rpc.Register: method %q has %d input parameters; needs exactly three\n",
          mname, mtype.NumIn()))
      continue
    }
    // First arg need not be a pointer.
    argType := make([]reflect.Type, 0, 6)
    pf_ := false
    for i := 2; i < inNums; i++ {
      ate := mtype.In(i)
      if !isExportedOrBuiltinType(ate) {
        msgErr = append(msgErr,
          fmt.Sprintf("rpc.Register: argument type of method %q is not exported: %q\n",
            mname, ate))
        pf_ = true
        break
      }

      if ate.Kind() == reflect.Interface {
        if i == inNums-1 {
          break
        } else {
          msgErr = append(msgErr,
            fmt.Sprintf("rpc.Register: argument io type of method %s must be lasted: %q\n",
              mname, ate))
          pf_ = true
          break
        }
      }

      argType = append(argType, ate)
    }

    if pf_ {
      continue
    }

    outNum := mtype.NumOut()
    replyType := make([]reflect.Type, 0)
    pf_ = false
    for i := 0; i < outNum; i++ {
      rte := mtype.Out(i)

      // Reply type must be exported.
      if !isExportedOrBuiltinType(rte) {
        msgErr = append(msgErr,
          fmt.Sprintf("rpc.Register: reply type of method %q is not exported: %q\n",
            mname, rte))
        pf_ = true
        break
      }

      if rte == typeOfError {
        if i < outNum-1 {
          // The last return type of the method must be error.
          msgErr = append(msgErr,
            fmt.Sprintf("rpc.Register: internal return type of method %q is %q, must not be error\n",
              mname, rte))

          pf_ = true
          break
        }
        break
      }

      replyType = append(replyType, rte)
    }

    if pf_ {
      continue
    }

    log.Debug("added method", mname)
    mname = strings.ToLower(mname)
    mets[mname] = &MethodType{Method: method, ArgType: argType, ReplyType: replyType}
  }

  return mets, msgErr
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
  for t.Kind() == reflect.Pointer {
    t = t.Elem()
  }
  // PkgPath will be non-empty even for an exported type,
  // so we need to check the type name as well.
  return token.IsExported(t.Name()) || t.PkgPath() == ""
}

// Precompute the reflect type for error.
var (
  typeOfError        = reflect.TypeFor[error]()
  typeOfReader       = reflect.TypeFor[io.Reader]()
  typeOfWriter       = reflect.TypeFor[io.Writer]()
  typeOfReaderWriter = reflect.TypeFor[io.ReadWriter]()
)
