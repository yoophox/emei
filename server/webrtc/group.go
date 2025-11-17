package webrtc

type otochan struct {
  wid   string
  peers [2]*peer
}

func (o *otochan) hSig(m *msg) error {
  switch m.Class {
  case SigClassOTOCallee:
    o.peers[1] = &peer{
      class: m.Class,
      typ:   m.Typ,
    }
    if m.Typ == SigTypeStmWebrtc && o.peers[0].typ == SigTypeStmWebrtc {
      return nil
    }
  }
  return nil
}

type otmroom struct {
  wid      string
  self     *peer
  watchers []*peer
}

type mtmroom struct {
  wid string
  pct []*peer
}

type lctroom struct {
  wid string
  tch *peer
  stu []*peer
}

type group interface {
  hSig(m *msg) error
  hData() // track and data channel
}
