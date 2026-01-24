package kube

import "github.com/yoophox/emei/cfg"

var (
  _etc      string
  _localDir string

  _svcCfgs    = svcCache{svcs: map[string]cfg.Config{}}
  _selfSvcCfg cfg.Config
)

const (
  CFG_RPC_PORT_PATH = "spec.ports.rpc"
  CFG_RPC_NET_PATH  = "metadata.lables.rpc-net"
)
