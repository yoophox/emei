package kube

import "github.com/yoophox/emei/cfg"

var (
  _etc      *bool
  _localDir *string

  _svcCfgs    = svcCache{svcs: map[string]cfg.Config{}}
  _selfSvcCfg cfg.Config
)

const (
  CFG_RPC_PORT_PATH   = "metadata.labels.server-port"
  CFG_RPC_NET_PATH    = "metadata.labels.rpc-net"
  CFG_ANNOTATIONS_PRE = "metadata.annotations."
)
