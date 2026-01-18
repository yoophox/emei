package cron

import (
  "time"

  "github.com/yolksys/emei/alg"
)

// ...
func AddJob(j *job, first *time.Time) {
  _elemChan <- &skiplistElem{j: j, t: first.Sub(_baseTime)}
}

// start ...
func schedule() {
  for {
    select {
    case m := <-_elemChan:
      _skiplist.Insert(m)
    case <-_ticker.C:
      _skiplist.GetSmallestNode().GetValue()
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
  _elemChan <- s
}

var (
  _baseTime time.Time
  _skiplist = alg.NewSkipList()
  _elemChan = make(chan *skiplistElem, 1000)
  _ticker   = time.NewTicker(time.Hour * 1000000)
)
