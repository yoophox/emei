package webrtc

import (
  "errors"
  "fmt"
  "net"

  "github.com/pion/stun"
)

func stunDo(buf []byte, addr net.Addr) (b []byte, err error) {
  // Parse the STUN message
  msg := new(stun.Message)
  msg.Raw = buf
  if err = msg.Decode(); err != nil {
    return
  }

  // Check if it's a Binding Request
  if msg.Type == stun.BindingRequest {
    // Create a Binding Response
    var response *stun.Message
    response, err = stun.Build(
      stun.TransactionID, stun.NewTransactionIDSetter(msg.TransactionID),
      stun.BindingSuccess,
      stun.XORMappedAddress{
        IP:   addr.(*net.UDPAddr).IP,
        Port: addr.(*net.UDPAddr).Port,
      },
    )
    if err == nil {
      b = response.Raw
    }

    return
  }

  err = fmt.Errorf("fail:stundo, reason: no supported msgtype, type:%d", msg.Type)
  return
}

// build ...
func buildErrRetv() ([]byte, error) {
  retv, err := stun.Build(
    stun.TransactionID,
    stun.BindingSuccess, stun.BindingError, stun.CodeBadRequest,
  )
  if err != nil {
    return nil, err
  }
  return retv.Raw, nil
}

var (
  _errDecode    = errors.New("decode")
  _errBuildRetv = errors.New("")
)
