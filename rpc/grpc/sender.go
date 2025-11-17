package grpc

import (
  "context"
  "crypto/tls"
  "encoding/json"
  "fmt"
  "net"
  "reflect"
  "time"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/rpc/call"
  "github.com/yolksys/emei/rpc/coder"
  "github.com/yolksys/emei/rpc/rpcabs"
  "github.com/quic-go/quic-go"
  rgrpc "google.golang.org/grpc"
)

type rvalue struct {
  r *Response
}

func (v *rvalue) RValues(rTyp ...reflect.Type) ([]reflect.Value, error) {
  if v.r.MsgH != nil && v.r.MsgH.Code != nil {
    return coder.DefaultValues(nil, rTyp...), fmt.Errorf("fail:grpc returned, reason:%s", *v.r.MsgH.Reason)
  }
  return coder.JsonStrArrToRvlaues(v.r.ResValues, rTyp...)
}
func (v *rvalue) Close() {}

// newGrpcClient ...
func newGrpcClient(svc *rpcabs.SvcInfo) error {
  tlsConfig := &tls.Config{
    InsecureSkipVerify: true,
  }
  qcfg := quic.Config{
    KeepAlivePeriod: 10 * time.Second,
  }
  creds := NewCredentials(tlsConfig)
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  qconn, err := quic.DialAddr(ctx, svc.Ip+":"+svc.Port, tlsConfig, &qcfg)
  if err != nil {
    return err
  }
  dialer := func(ctx context.Context, _ string) (net.Conn, error) {
    stream, err := qconn.OpenStream()
    if err != nil {
      return nil, err
    }

    return &QuicConn{qconn, stream}, nil
  }

  c, err := rgrpc.NewClient(svc.Ip, rgrpc.WithContextDialer(dialer), rgrpc.WithTransportCredentials(creds))
  if err != nil {
    return err
  }
  svc.Client = NewGrpcClient(c)

  return nil
}

type sender struct {
  rpcabs.SndImp
}

func (s *sender) Send(e env.Env, svc, met string, args ...any) call.Response {
  defer e.Return()

  _svc, err := s.GetSvc(svc, newGrpcClient)
  e.AssertErr(err)
  c := _svc.Client.(GrpcClient)
  h := e.GetMsgHeader()
  params := []string{}
  for _, value := range args {
    _p, err := json.Marshal(value)
    e.AssertErr(err)
    params = append(params, string(_p))
  }
  req := &Request{
    Met: met,
    MsgH: &Header{
      Mid: h.Mid,
      Jwt: h.Jwt,
      Sid: h.Sid,
    },
    ReqParams: params,
  }
  r, err := c.Route(context.Background(), req)
  e.AssertErr(err)
  return &rvalue{r}
}

// newSender ...
func newSender() *sender {
  s := &sender{}
  s.For = "grpc"
  return s
}
