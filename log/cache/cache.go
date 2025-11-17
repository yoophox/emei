package cache

import (
  "sync"
)

type (
  Content []byte
)

type LogRecord struct {
  TraceId string
  Attris  map[string]string
  Buf     Content
}

func Put(t Content) {
  // Proper usage of a sync.Pool requires each entry to have approximately
  // the same memory cost. To obtain this property when the stored type
  // contains a variably-sized buffer, we add a hard limit on the maximum buffer
  // to place back in the pool.
  //
  // See https://golang.org/issue/23199
  const maxSize = 1 << 16 // 64KiB
  if cap(t) > maxSize {
    return
  }
  _cachePool.Put(t)
}

func Get() Content {
  return _cachePool.Get().(Content)
}

// GetLogRecord ...
func GetLogRecord() *LogRecord {
  return _logRecPool.Get().(*LogRecord)
}

// ReleaseLogRecord ...
func ReleaseLogRecord(r *LogRecord) {
  _logRecPool.Put(r)
}

var _cachePool = &sync.Pool{
  New: func() any {
    return make(Content, 0, 1000)
  },
}

var _logRecPool = sync.Pool{
  New: func() any {
    return &LogRecord{
      Attris: map[string]string{},
    }
  },
}
