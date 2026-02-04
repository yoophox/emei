package uuid

import (
  "fmt"
  "time"
)

// Local ...
func Local() int64 {
  uu_, err := _localIdGener.NextID()
  if err != nil {
    panic(err)
  }
  return uu_
}

// LocalStr ...
func LocalHex() string {
  uuid := Local()
  return fmt.Sprintf("%016X", uuid)
}

// Sevice ...
func Sevice() int64 {
  // --uuid.machineid > 0
  // stateful in kubernates
  // @@valkey service exist
  // or @@uuid exist
  if _machineId > 100 {
    return Local()
  }
  return UUID()
}

// Service ...
func ServiceHex() string {
  return ""
}

// uuid ...
func UUID() int64 {
  // @@uuid: gob of number
  // @@uuid.get(190)

  if !_uuidWorking {
    return 0
  }
  timer := time.NewTimer(time.Second)
  defer timer.Stop()
  select {
  case <-timer.C:
    return 0
  case uuid := <-_uuidCh:
    return uuid
  }
}

// String ...
func Hex() string {
  return ""
}
