package cron

import "time"

func init() {
  _baseTime = time.Now()
  go schedule()
}
