package core

import (
  "path"
  "runtime"
  "strconv"
  "time"

  "github.com/yolksys/emei/log/cache"
)

func caller(l *logger) {
  if l.cSkip < 0 {
    return
  }
  pc, file, line, ok := runtime.Caller(l.cSkip)
  if !ok {
    return
  }
  l.buf = serialize(l.buf, l.cFName, callerMarshalFunc(pc, file, line))
}

func timeField(l *logger) {
  if l.tFmt == "" {
    return
  }

  l.buf = serialize(l.buf, l.tFName, time.Now().Format(l.tFmt))
}

func write(l *logger) {
  buf := l.buf
  rc_ := l.LogRecord
  l.buf = cache.Get()[:0]
  l.buf = enc.AppendBeginMarker(l.buf)
  l.LogRecord = cache.GetLogRecord()
  buf = enc.AppendEndMarker(buf)
  rc_.Buf = buf

  for _, b := range bcd {
    b.Write(rc_)
  }
}

func (l *logger) doLog(lev string, m ...any) Log {
  // if l.cacheMode {
  //   l.buf = append(l.buf, "\n~"...)
  // }
  l.buf = serialize(l.buf, l.lFName, lev)
  caller(l)
  timeField(l)
  l.buf = serialize(l.buf, m...)
  if !l.isCacheMod {
    write(l)
  } else {
    l.buf = append(l.buf, "\n ~"...)
  }
  return l
}

var (
  // CallerMarshalFunc allows customization of global caller marshaling
  callerMarshalFunc = func(pc uintptr, file string, line int) string {
    fName := runtime.FuncForPC(pc).Name()
    return path.Base(file) + ":" + path.Base(fName) + ":" + strconv.Itoa(line)
  }

  timeDefaultFmt = "2006-01-02 15:04:05.999999999 -0700 MST"
)
