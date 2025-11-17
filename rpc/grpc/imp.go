package grpc

import (
  "context"
  "crypto/tls"
  "errors"
  "fmt"
  "reflect"
  "strings"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/rpc/coder"
  "github.com/yolksys/emei/rpc/rpcabs"
  "github.com/quic-go/quic-go"
  rgrpc "google.golang.org/grpc"
)

type recver struct {
  stm *stream
  i   *impl
}

func (i *recver) Route(_ context.Context, r *Request) (*Response, error) {
  tlk := &Response{}
  e := env.New("grpc", r.Met, newEnc(tlk), newDec(r))
  defer e.Finish()

  _path := strings.Split(r.Met, ".")
  if len(_path) != 2 {
    e.AssertErr(errors.New("fail:met,met:" + r.Met))
  }

  recv, ok := i.i.Recvs[_path[0]]
  if !ok {
    e.AssertErr(errors.New("fail: no recver,name:" + _path[0]))
  }

  met, ok := recv.Mets[_path[1]]
  if !ok {
    e.AssertErr(errors.New("fail: no method,name:" + _path[1]))
  }
  parms := coder.ParseParam(e, met, recv)
  e.Assert()
  e.PrintParams(parms...)
  if i.stm != nil {
    parms = append(parms, reflect.ValueOf(i.stm))
  }

  f := met.Method.Func
  rtv := f.Call(parms)
  if i.stm == nil {
    e.SetReV(rtv)
    return tlk, nil
  }

  if e.CheckErr(rtv...) {
    return nil, e.Err()
  }

  return nil, nil
}

func (rpc *recver) RWStream(stm Grpc_RWStreamServer) error {
  m, err := stm.Recv()
  if err != nil {
    return err
  }

  if m.Req == nil {
    return fmt.Errorf("fail:grpc rwstream, reason:m.req is nil")
  }
  rpc.stm = &stream{c: stm}
  _, err = rpc.Route(nil, m.Req)
  return err
}

func (rpc *recver) mustEmbedUnimplementedGrpcServer() {}

type impl struct {
  rpcabs.RpcImp
}

// NewRpc ...
func NewRpc() rpcabs.RPC {
  s := &impl{}
  s.RpcImp.Name = "grpc"
  return s
}

func (i *impl) Name() string {
  return i.RpcImp.Name
}

func (i *impl) Start() error {
  qcfg := quic.Config{Allow0RTT: true}
  tlscfg := tls.Config{
    Certificates: []tls.Certificate{i.Certificate},
  }
  ln_, err := quic.ListenAddr(":"+i.Port, &tlscfg, &qcfg)
  if err != nil {
    return err
  }
  log.Event("grpc listen port", i.Port)

  go func() {
    // ctx, c := context.WithTimeout(context.Background(), 3*time.Second)
    // defer c()
    for {
      conn, err := ln_.Accept(context.Background())
      if err != nil {
        i.Err <- err
        break
      }
      gln := Listen(ln_, conn)
      go func() {
        for {
          log.Event("msg", "start a new grpc server")
          server := rgrpc.NewServer()
          RegisterGrpcServer(server, &recver{i: i})
          err = server.Serve(gln)
          log.Event("msg", "fail: accept stream", "error", err)
          break
        }
      }()
    }
    ln_.Close()
  }()

  return nil
}

func (i *impl) Close() {
}
