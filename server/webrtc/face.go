package webrtc

// signal class from 1------------------------------
const (
  // req attrs: uid=userid
  // retv attrs: optinal code, optional reason
  SigClassLogin byte = iota + 1
  // req attrs: target=userid
  // retv attrs: optional code, optional reason, wid
  SigClassOTOCall // one to one
  // req attrs: wid
  // retv attrs: optional code, optional reason
  SigClassOTOCallee // one to one
  // req attrs: optional passd
  // retv attrs: optional code, optional reason, wid
  SigClassAllocOTM // one to multiple
  SigClassAllocMTM // multi to multi
  // Lecture:Students can only see the teacher
  // Or can see vidio and hear voice from the student who is currently speaking
  // teachers can see all students
  SigClassAllocLCT
  // req attrs: passd, wid
  // retv attrs: optional code, optional reason
  SigClassAllocCli = 15 // Participants
)

const (
  // req attrs: wid
  SigClassCfg byte = iota + 16
)

const (
  SigClassNotify   byte = 251
  SigClassResponse byte = 255
)

// signal type from 1---------------------------------
const (
  SigTypeStmWebrtc byte = iota + 1
  SigTypeStmRTP
  SigTypeStmRTMP

  SigStmTypeCustomTCP  = 13
  SigStmTypeCustomUDP  = 14
  SigStmTypeCustomQUIC = 15
)

const (
  SigTypeCfgSDP byte = iota + 16
  SigTypeCfgCadidate
)

// attrs: caller="name"
const SigTypeNtyCaller byte = iota + 235
