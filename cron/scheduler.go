package cron

import (
  "time"

  "github.com/yolksys/emei/alg"
)

// ...
func AddJob(j *job, first *time.Time) *skiplistElem {
  e := &skiplistElem{j: j, t: first.Sub(_baseTime)}
  _elemChan <- e
  return e
}

// start ...
func schedule() {
  for {
    select {
    case m := <-_elemChan:
      _skiplist.Insert(m)
    case m := <-_removeElemChan:
      _skiplist.Delete(m)
    case <-_ticker.C:
      for {
        m := _skiplist.GetSmallestNode().GetValue()
        if m == nil {
          _ticker.Reset(time.Hour * 1000000)
          break
        }
        s := m.(*skiplistElem)
        dur := s.t - time.Since(_baseTime)
        if dur <= 0 {
          go s.doJob()
          _skiplist.Delete(m)
          continue
        }

        _ticker.Reset(dur)
      }
    }
  }
}

// skip-list

func init() {
  _baseTime = time.Now()
  schedule()
}

type skiplistElem struct {
  j *job
  t time.Duration // sub(basetime)
}

func (s *skiplistElem) ExtractKey() float64 {
  return float64(s.t)
}

func (s *skiplistElem) String() string {
  return ""
}

func (s *skiplistElem) doJob() {
  t := s.j.Run()
  if t == nil {
    return
  }

  s.t = t.Sub(_baseTime)
  _elemChan <- s
}

var (
  _baseTime       time.Time
  _skiplist       = alg.NewSkipList()
  _elemChan       = make(chan *skiplistElem, 1000)
  _removeElemChan = make(chan *skiplistElem, 1000)
  _ticker         = time.NewTicker(time.Hour * 1000000)
)
