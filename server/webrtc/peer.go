package webrtc

import (
  "github.com/pion/webrtc/v4"
)

type peer struct {
  class byte // otm, mtm lct...
  typ   byte // stream type: "rtpc/weprtc"
  name  string
  sigw  chan *msg // write ice to peer
  stmr  stmReader
  stmw  stmWriter
}

func (p *peer) Init() error {
  switch p.typ {
  case SigTypeStmWebrtc:
  }

  return nil
}

type stmData struct {
  stmType byte
  video   []byte
  audio   []byte
}

type stmReader interface {
  Read() (*stmData, error)
}

type stmWriter interface {
  Writte(d *stmData) error
}

// newPeer ...
func newPeer(sigclass, stmTyp byte, name string) (*peer, error) {
  p := &peer{typ: stmTyp, class: sigclass, name: name}
  switch sigclass {
  case SigClassOTOCall:
  case SigClassOTOCallee:
  case SigClassAllocOTM:
  case SigClassAllocMTM:
  case SigClassAllocLCT:
  }

  return p, nil
}

type webrtcReader struct {
  tracks []*webrtc.TrackRemote
}

type webrtcWriter struct {
  tracks []*webrtc.TrackLocal
}
