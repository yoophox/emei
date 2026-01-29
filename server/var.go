package svr

import (
  "io"
  "net/http"
  "reflect"

  "github.com/yoophox/emei/env"
)

const (
  SERVER_FOR_WEB SvrFor = "web"
  SERVER_FOR_RPC SvrFor = "rpc"
)

const (
  RPC_NET_QUIC netTx = iota
  RPC_NET_TCP
)

var (
  RPC_JWT_COOKIES_NAME = "LPKJNHHEjhdaferHVNOGCCFDSXsdjewjsdjWLKZMEXS-jwt"
  RPC_VER_COOKIES_NAME = "AFEdfuywbndskjewhccjsdjhLWTTCBMMAWFVNJDENCD-ver"
  // traceid
  RPC_TID_COOKIES_NAME = "AFEdfuywbndskjsdfeltsdjhLWTTCBMMAWFVNJDENCD-tid"
)

const (
  RPC_ALP_GOB   tlsAlpnTx = "gob"
  RPC_ALP_GRPC  tlsAlpnTx = "grpc"
  RPC_ALP_JSON  tlsAlpnTx = "json"
  RPC_ALP_HTTP2 tlsAlpnTx = "h2"
  RPC_ALP_HTTP3 tlsAlpnTx = "h3"
)

var (
  _rpcRecvs = map[string]*rcvrTx{}
  _webRecvs = map[string]*rcvrTx{}
  //_wg       *sync.WaitGroup
)

// Precompute the reflect type for error.
var (
  typeOfError        = reflect.TypeFor[error]()
  typeOfReader       = reflect.TypeFor[io.Reader]()
  typeOfWriter       = reflect.TypeFor[io.Writer]()
  typeOfReaderWriter = reflect.TypeFor[io.ReadWriter]()
  typeOfWebsock      = reflect.TypeFor[WebSock]()
  typeOfWebRes       = reflect.TypeFor[WebResponse]()
  typeOfRequest      = reflect.TypeFor[*http.Request]()
  typeOfEnv          = reflect.TypeFor[env.Env]()
)

var _rootEnv env.Env //= env.New(nil)

const (
  _RPC_TIMEOUT = 5 // second
)
