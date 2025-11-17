package nrpc

import (
  "context"
  "crypto/tls"
  "encoding/json"
  "fmt"
  "io"
  "reflect"
  "runtime"
  "sync"
  "sync/atomic"
  "time"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/call"
  "github.com/yolksys/emei/rpc/coder"
  "github.com/yolksys/emei/rpc/rpcabs"
  "github.com/quic-go/quic-go"
)

// for client...
func newClient(svc *rpcabs.SvcInfo) error {
  ctx, c := context.WithTimeout(context.Background(), time.Second*3)
  defer c()
  qcfg := quic.Config{Allow0RTT: true}
  tlscfg := pki.NewClientTlsConfig() // tls.Config{InsecureSkipVerify: true}
  conn, err := quic.DialAddr(ctx, svc.Ip+":"+svc.Port, tlscfg, &qcfg)
  if err != nil {
    return err
  }

  newRpcClient := func(isStm bool) any {
    stm, err := conn.OpenStream()
    if err != nil {
      return err
    }

    clr := newNrpcCaller(stm)
    clr.isStm = isStm
    runtime.AddCleanup(clr, rpcclientCloser, clr.codec.Close)
    return clr
  }

  cli := &client{
    c: conn,
    callPool: &sync.Pool{
      New: func() any {
        return newRpcClient(false)
      },
    },
    stmCallPool: &sync.Pool{
      New: func() any {
        return newRpcClient(true)
      },
    },
  }
  cli.usedPool = cli.callPool
  svc.Client = cli
  return nil
}

// rpcclientCloser ...
func rpcclientCloser(c func() error) {
  // fmt.Println("nrpc-funcs.go-rpcclientCloser")
  c()
}

// for sender...----------------------------------------------------
func newSender() *sender {
  return &sender{SndImp: rpcabs.SndImp{For: "nrpc"}}
}

func (s *sender) Send(e env.Env, svc, met string, args ...any) call.Response {
  defer e.Return()

  si_, err := s.GetSvc(svc, newClient)
  e.AssertErr(err)
  c := si_.Client.(*client).usedPool.Get()
  del := func() {
    s.Delete(svc)
  }
  switch m := c.(type) {
  case error:
    e.AssertErr(m, del)
  case *nrpcCaller:
    m.req.MsgH = e.GetMsgHeader()
    err = m.Call(met, args...)
    e.AssertErr(err, del)
    var resp response
    resp.clr = m
    if !m.isStm {
      si_.Client.(*client).usedPool.Put(m)
    }
    return &resp
  }
  return nil
}

func (s *sender) SendWithStream(e env.Env, svc, met string, args ...any) io.ReadWriteCloser {
  defer e.Return()

  si_, err := s.GetSvc(svc, newClient)
  e.AssertErr(err)

  c := si_.Client.(*client)
  c.usedPool = c.stmCallPool
  defer func() {
    c.usedPool = c.callPool
  }()
  t := s.Send(e, svc, met, args...)
  e.Assert()
  tt_ := t.(*response)
  if tt_.clr.retv.MsgH != nil && tt_.clr.retv.MsgH.Code != 0 {
    c.stmCallPool.Put(tt_.clr)
    e.AssertErr(fmt.Errorf("fail:SendWithStream, code:1, reason:%s", tt_.clr.retv.MsgH.Reason))
  }

  close := func() {
    c.stmCallPool.Put(tt_.clr)
  }

  tt_.clr.stm.close = close
  tt_.clr.stm.s = tt_.clr.Stream
  return &tt_.clr.stm
}

func (w *response) RValues(rTyp ...reflect.Type) ([]reflect.Value, error) {
  mh_ := w.clr.retv.MsgH
  if mh_ != nil && mh_.Code != 0 {
    return coder.DefaultValues(nil, rTyp...), fmt.Errorf("fail:nrpc, reason:%s", mh_.Reason)
  }
  return coder.JsonStrArrToRvlaues(w.clr.retv.Data, rTyp...)
}

func (w *response) Close() {
}

// for rpc ------------------------------------------------------
// NewRpc ...
func NewRpc() rpcabs.RPC {
  s := &nrpc{}
  s.RpcImp.Name = "nrpc"
  _nrpc = s
  return s
}

func (s *nrpc) Name() string {
  return s.RpcImp.Name
}

func (s *nrpc) Start() error {
  qcfg := quic.Config{Allow0RTT: true}
  tlscfg := tls.Config{
    Certificates: []tls.Certificate{s.Certificate},
  }
  log.Event("nrpc", "listen", "port", s.Port)
  ln_, err := quic.ListenAddr(":"+s.Port, &tlscfg, &qcfg)
  if err != nil {
    return err
  }

  go func() {
    // ctx, c := context.WithTimeout(context.Background(), 0*time.Second)
    // defer c()
    for {
      conn, err := ln_.Accept(context.Background())
      if err != nil {
        s.Err <- err
        break
      }

      go func() {
        connStmCnt := int32(0)
        var stmCntTime *time.Time
        for {
          stream, err := conn.AcceptStream(context.Background())
          if err != nil {
            log.Event("msg", "fail: nrpc accept stream", "error", err)
            break
          }

          log.Event("nrpc", "start a server")
          go func() {
            svr := newServer(stream)
            defer svr.Close()
            for {
              atomic.AddInt32(&connStmCnt, 1)
              if connStmCnt > 10 {
                if stmCntTime == nil {
                  t := time.Now()
                  stmCntTime = &t
                }
              } else {
                stmCntTime = nil
              }
              err := svr.serve(&connStmCnt)
              if connStmCnt > 10 && time.Now().Sub(*stmCntTime) > time.Hour {
                return
              }
              if err != nil {
                log.Event("nrpc", "stop a server")
                return
              }
            }
          }()
        }
      }()
    }
    ln_.Close()
  }()

  return nil
}

func (s *nrpc) Close() {
}

// for caller----------------------------------------------------------------
func (c *nrpcCaller) Call(met string, args ...any) error {
  // rreq := rpc.Request{ServiceMethod: met}
  c.Request.ServiceMethod = met
  c.req.Data = c.req.Data[:0]
  if c.isStm {
    c.req.IsStm = &c.isStm
  } else {
    c.req.IsStm = nil
  }
  for _, value := range args {
    str, _ := json.Marshal(value)
    c.req.Data = append(c.req.Data, string(str))
  }

  err := c.codec.WriteRequest(&c.Request, &c.req)
  if err != nil {
    return fmt.Errorf("fail:wirte request, reason:%s}", err.Error())
  }

  err = c.codec.ReadResponseHeader(&c.Response)
  if err != nil {
    return fmt.Errorf("fail: read response, reason:%s", err.Error())
  }

  // c.retv.Data = c.req.Data[:0]
  c.retv.MsgH.Code = 0
  c.retv.MsgH.Reason = ""
  err = c.codec.ReadResponseBody(&c.retv)
  if err != nil {
    return fmt.Errorf("fail:read response body, reason:%s", err.Error())
  }
  return nil
}

func (c *nrpcCaller) Close() error {
  return c.Stream.Close()
}

// newNrpcCaller ...
func newNrpcCaller(q quic.Stream) *nrpcCaller {
  nc_ := &nrpcCaller{
    Stream: q,
    codec:  newCliGobCodec(q),
    isStm:  false,
  }

  nc_.req.Data = []string{}
  nc_.retv.MsgH = &env.Tjatse{}
  return nc_
}

// for server-----------------------------------------------------
func (s *server) Close() {
  s.Stream.Close()
}

func (s *server) stmResonse() error {
  return s.WriteResponse(&s.Response, &s.res)
}

func (s *server) served() error {
  if s.req.IsStm != nil && *s.req.IsStm && s.isStmResed {
    return s.stm.Close()
  }

  return s.gobServerCodec.WriteResponse(&s.Response, &s.res)
}

// newServer ...
func newServer(q quic.Stream) *server {
  s := &server{
    Stream:         q,
    gobServerCodec: newSvrGobCodec(q),
    res:            NRpcData{Data: []string{}},
  }
  s.stm.s = q
  return s
}
