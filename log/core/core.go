package core

import (
  "context"
  "sync"

  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/log/backend"
  "github.com/yolksys/emei/log/cache"
  "github.com/yolksys/emei/utils"
)

// Level defines log levels.
type Level int8

const (
  // DebugLevel defines debug log level.
  DebugLevel Level = iota
  // InfoLevel defines info log level.
  InfoLevel
  // WarnLevel defines warn log level.
  WarnLevel
  // ErrorLevel defines error log level.
  EventLevel
  // ErrorLevel defines error log level.
  ErrorLevel
  // FatalLevel defines fatal log level.
  FatalLevel
  // PanicLevel defines panic log level.
  PanicLevel
  // NoLevel defines an absent log level.
  NoLevel
  // Disabled disables the logger.
  Disabled

  // TraceLevel defines trace log level.
  TraceLevel Level = -1
  // Values less than TraceLevel are handled as numbers.
)

type Log interface {
  Debug(...interface{}) Log
  Info(...interface{}) Log
  Warn(...interface{}) Log
  Event(...interface{}) Log
  Error(...interface{}) Log
  Fatal(...interface{}) Log
  Log(...interface{}) Log
  Flush() Log

  Level(Level) Log
  Prefix(...interface{}) Log
  CallerSkip(int) Log    // if < 0 disable caller field
  TimeFmt(string) Log    // = "" disable time field
  TFieldName(string) Log // time name
  CFieldName(string) Log // Caller name
  LFieldName(string) Log // Level name
  // ConfigCacheMode(c bool) Log
  SetTraceId(id string) Log
  AddAttri(key, value string) Log
  // Suffix(...interface{}) Log
}

func New(ctx context.Context, opts ...optfunc) Log {
  l := loggerPool.Get().(*logger)
  l.buf = cache.Get()
  l.buf = l.buf[:0]
  l.LogRecord = cache.GetLogRecord()
  resetOpt(&l.option)
  for _, value := range opts {
    value(&l.option)
  }

  l.buf = enc.AppendBeginMarker(l.buf)
  if l.isCacheMod {
    l.buf = append(l.buf, '~')
  }

  l.level = level
  l.ctx = ctx
  l.tFmt = timeDefaultFmt
  l.tFName = "T"
  l.cFName = "C"
  l.lFName = "L"
  l.cSkip = baseCallerSkip

  return l
}

type logger struct {
  buf cache.Content

  level  Level
  ctx    context.Context
  tFmt   string
  tFName string // time field name
  cFName string // caller field name
  lFName string // level field name
  cSkip  int    // caller skip

  *cache.LogRecord

  option
}

func (l *logger) Debug(m ...interface{}) Log {
  if l.level > DebugLevel {
    return l
  }

  return l.doLog("debug", m...)
}

func (l *logger) Info(m ...interface{}) Log {
  if l.level > InfoLevel {
    return l
  }

  return l.doLog("info", m...)
}

func (l *logger) Warn(m ...interface{}) Log {
  if l.level > WarnLevel {
    return l
  }

  return l.doLog("warn", m...)
}

func (l *logger) Event(m ...interface{}) Log {
  if l.level > EventLevel {
    return l
  }

  return l.doLog("event", m...)
}

func (l *logger) Error(m ...interface{}) Log {
  if l.level > ErrorLevel {
    return l
  }

  return l.doLog("error", m...)
}

func (l *logger) Fatal(m ...interface{}) Log {
  if l.level > FatalLevel {
    return l
  }

  return l.doLog("fatal", m...)
}

func (l *logger) Log(m ...interface{}) Log {
  caller(l)
  timeField(l)
  l.buf = serialize(l.buf, m...)
  if !l.isCacheMod {
    write(l)
  }
  return l
}

func (l *logger) Flush() Log {
  if l.isCacheMod {
    write(l)
  }
  return l
}

func (l *logger) Prefix(m ...interface{}) Log {
  l.buf = serialize(l.buf, m...)
  return l
}

func (l *logger) Level(lev Level) Log {
  l.level = lev
  return l
}

func (l *logger) CallerSkip(s int) Log {
  if s < 0 {
    l.cSkip = -1
    return l
  }

  l.cSkip = baseCallerSkip + s
  return l
}

func (l *logger) TimeFmt(f string) Log {
  l.tFmt = f
  return l
}

func (l *logger) TFieldName(n string) Log {
  l.tFName = n
  return l
}

func (l *logger) CFieldName(n string) Log {
  l.cFName = n
  return l
}

func (l *logger) LFieldName(n string) Log {
  l.lFName = n
  return l
}

func (l *logger) ConfigCacheMode(c bool) Log {
  l.isCacheMod = c
  return l
}

func (l *logger) SetTraceId(id string) Log {
  l.TraceId = id
  return l
}

func (l *logger) AddAttri(key, value string) Log {
  l.Attris[key] = value
  return l
}

func Release(_l Log) {
  l := _l.(*logger)
  if l.buf != nil {
    cache.Put(l.buf)
    l.buf = nil
  }

  if l.LogRecord != nil {
    cache.ReleaseLogRecord(l.LogRecord)
    l.LogRecord = nil
  }

  loggerPool.Put(l)
}

var (
  loggerPool = &sync.Pool{
    New: func() any {
      return &logger{}
    },
  }

  level Level = 0

  bcd []backend.Backend = make([]backend.Backend, 0)

  baseCallerSkip = 3
)

func init() {
  // new backends due to cfg
  var bcds []string
  err := cfg.GetCfgItem("logger.level", &level)
  utils.AssertErr(err)
  err = cfg.GetCfgItem("logger.backends", &bcds)
  utils.AssertErr(err)
  for _, v := range bcds {
    f := backend.Get(v)
    if f == nil {
      continue
    }
    bcd = append(bcd, f(context.Background()))
  }
}
