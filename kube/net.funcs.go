package kube

import "github.com/yoophox/emei/utils"

func (n *Net) GetPortByName(name string) string {
  // c := getSvcCfg()
  nm_ := ""
  _selfSvcCfg.Scan(CFG_ANNOTATIONS_PRE+n.svc+"-"+name, &nm_)
  if nm_ != "" {
    name = nm_
  }

  p, _ := n.Ports[name]
  return p.Port
}

func (s *Server) Addr(portName string) string {
  port := s.GetPortByName(portName)
  return utils.CompriseAddr(s.IP, port)
}
