package svr

import "net/http"

// serveHttp ...
func serveHttp(w http.ResponseWriter, r *http.Request) {
  cntTyp := r.Header.Get("content-type")
  if r.Header.Get("upgrade") == "websocket" {
  } else if cntTyp == "multipart/form-data" {
    err := r.ParseMultipartForm(1000 * 1000 * 100)
    if err != nil {
      w.Write([]byte(err.Error()))
      w.WriteHeader(http.StatusBadRequest)
      return
    }
    // stream = (*upFiler)(r.MultipartForm)
  } else if cntTyp == "application/json" {
  } else {
    w.Write([]byte("wrong content type:" + cntTyp))
    w.WriteHeader(http.StatusBadRequest)
  }
}
