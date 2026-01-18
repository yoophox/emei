package errors

import "github.com/yolksys/emei/errs"

const (
  ERR_ID_RPC_QUIC_LISTEN      errs.ErrId = "err.rpc.quic.listen"
  ERR_ID_RPC_TCP_LISTEN       errs.ErrId = "err.rpc.tcp.listen"
  ERR_ID_RPC_TCP_DIAL         errs.ErrId = "err.rpc.tcp.dial"
  ERR_ID_RPC_NEW_SESSION      errs.ErrId = "err.rpc.new.session"
  ERR_ID_RPC_QUIC_DIAL_EARLY  errs.ErrId = "err.rpc.dial.early"
  ERR_ID_RPC_QUIC_OPEN_STREAM errs.ErrId = "err.rpc.open.stream"
  ERR_ID_RPC_DECODE           errs.ErrId = "err.rpc.decode"
  ERR_ID_RPC_ENCODE           errs.ErrId = "err.rpc.encode"
  ERR_ID_RPC_CALL_SEND        errs.ErrId = "err.rpc.call.send"
  ERR_ID_RPC_DECODE_TJAX      errs.ErrId = "err.rpc.decode.tjax"
  ERR_ID_RPC_DECODE_CALLINFO  errs.ErrId = "err.rpc.decode.callinfo"
  ERR_ID_RPC_CALLINFO_LEN     errs.ErrId = "err.rpc.callinfo.len"
  ERR_ID_RPC_CALLINFO_RCVR    errs.ErrId = "err.rpc.callinfo.rcvr"
  ERR_ID_RPC_CALLINFO_METH    errs.ErrId = "err.rpc.callinfo.meth"
  ERR_ID_RPC_DECODE_PARAMS    errs.ErrId = "err.rpc.decode.params"
)
