package svr

import "net/http"

// newWebResponse ...
func newWebResponse() *webResponsImpl {
  return &webResponsImpl{
    headers: map[string]string{},
    cookies: []*http.Cookie{},
    content: [][]byte{},
  }
}

func (r *webResponsImpl) AddHeader(k, v string) {
  r.headers[k] = v
}

func (r *webResponsImpl) AddCookie(c *http.Cookie) {
  r.cookies = append(r.cookies, c)
}

func (r *webResponsImpl) Write(b []byte) {
  r.content = append(r.content, b)
}
