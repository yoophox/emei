package web

import (
  "io"
  "mime/multipart"
  "sync"
)

type UpFile interface {
  //  Values() map[string][]string
  Files() map[string][]*multipart.FileHeader
}

type upFiler multipart.Form

func (u *upFiler) Values() map[string][]string {
  return (*multipart.Form)(u).Value
}

func (u *upFiler) Files() map[string][]*multipart.FileHeader {
  return (*multipart.Form)(u).File
}

type DnFiler interface {
  Path() string
  io.Writer
}

// newDnFiler ...
func newDnFiler(w io.Writer, p string) *dnFiler {
  d := _dnFilePool.Get().(*dnFiler)
  d.w = w
  d.path = p
  return d
}

type dnFiler struct {
  path string
  w    io.Writer
}

func (d *dnFiler) Path() string {
  return d.path
}

func (d *dnFiler) Write(buf []byte) (int, error) {
  return d.w.Write(buf)
}

var _dnFilePool = sync.Pool{
  New: func() any {
    return &dnFiler{}
  },
}
