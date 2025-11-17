package log

import (
  "context"
  "sync"

  "github.com/yolksys/emei/log/core"
)

type Logger = core.Log

var (
  New           = core.New
  Release       = core.Release
  WithCacheMode = core.WithCacheMode()
)

// DebugLevel defines debug log level.
var DebugLevel = core.DebugLevel

// InfoLevel defines info log level.
var InfoLevel = core.InfoLevel

// WarnLevel defines warn log level.
var WarnLevel = core.WarnLevel

// ErrorLevel defines error log level.
var EventLevel = core.EventLevel

// ErrorLevel defines error log level.
var ErrorLevel = core.ErrorLevel

// FatalLevel defines fatal log level.
var FatalLevel = core.FatalLevel

func Debug(f ...interface{}) Logger {
  mu.Lock()
  defer mu.Unlock()
  globalL.Debug(f...)
  return globalL
}

func Info(f ...interface{}) Logger {
  mu.Lock()
  defer mu.Unlock()
  globalL.Info(f...)
  return globalL
}

func Warn(f ...interface{}) Logger {
  mu.Lock()
  defer mu.Unlock()
  globalL.Warn(f...)
  return globalL
}

func Event(f ...interface{}) Logger {
  mu.Lock()
  defer mu.Unlock()
  globalL.Event(f...)
  return globalL
}

func Error(f ...interface{}) Logger {
  mu.Lock()
  defer mu.Unlock()
  globalL.Error(f...)
  return globalL
}

func Fatal(f ...interface{}) Logger {
  mu.Lock()
  defer mu.Unlock()
  globalL.Fatal(f...)
  return globalL
}

func Log(f ...interface{}) Logger {
  mu.Lock()
  defer mu.Unlock()
  globalL.Log(f...)
  return globalL
}

var (
  mu      = sync.Mutex{}
  globalL = New(context.Background()).CallerSkip(1)
)
