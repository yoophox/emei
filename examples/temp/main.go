package main

import (
  "github.com/yolksys/emei"
)

func main() {
  // used services by this app
  _ = []string{"@@service0"}

  emei.RegRecver(_sssss)
  emei.Serve()
}

func init() {
  emei.AssertCmd()
}
