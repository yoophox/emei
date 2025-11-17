package resolver

import (
  "context"
  "fmt"
  "net"
  "net/http"
  "strings"
  "time"

  "github.com/yolksys/emei/etc"
  "github.com/yolksys/emei/utils"
)

type conn struct {
  network, addr  string
  rwdl, rdl, wdl *time.Time
  rMsg           *DNSMessage
  err            error
}

func (c *conn) Read(buf []byte) (int, error) {
  if c.err != nil {
    return 0, c.err
  }

  c.rMsg.Header.ANCount = uint16(len(c.rMsg.Answers))
  c.rMsg.Header.ARCount = uint16(len(c.rMsg.AdditionalRRs))
  c.rMsg.Header.NSCount = uint16(len(c.rMsg.AuthorityRRs))
  c.rMsg.Header.QDCount = uint16(len(c.rMsg.Questions))
  // buf = buf[:0]
  b := c.rMsg.ToBytes(buf)
  return len(b), nil
}

func (c *conn) Write(b []byte) (n int, err error) {
  c.getRetValue(b)
  return len(b), nil
}

func (c *conn) Close() error {
  return nil
}

func (c *conn) LocalAddr() net.Addr {
  return nil
}

func (c *conn) RemoteAddr() net.Addr {
  return nil
}

func (c *conn) SetDeadline(t time.Time) error {
  c.rwdl = &t
  return nil
}

func (c *conn) SetReadDeadline(t time.Time) error {
  c.rdl = &t
  return nil
}

func (c *conn) SetWriteDeadline(t time.Time) error {
  c.wdl = &t
  return nil
}

func (c *conn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
  return 0, nil, nil
}

func (c *conn) WriteToa(p []byte, addr net.Addr) (n int, err error) {
  return 0, nil
}

func (c *conn) query(p []byte) []byte {
  for _, v := range _nameserver {
    if v == "etc" {
      ip, err := etc.GetSvcIp(c.rMsg.Questions[0].Name)
      if err != nil {
        continue
      }

      typ := TypeAAAA
      if utils.IsIpv4(ip) {
        typ = TypeA
      }
      c.rMsg.Answers = append(c.rMsg.Answers, &RR{
        Name:        c.rMsg.Questions[0].Name,
        Type:        typ,
        Class:       c.rMsg.Questions[0].QClass,
        RDLength:    (uint16)(len(ip)),
        RData:       []byte(ip),
        RDataParsed: ip,
      })
      return []byte{1}
    }

    d := net.Dialer{}
    c_, err := d.Dial(c.network, v)
    if err != nil {
      c.err = err
      return nil
    }

    if c.rwdl != nil {
      c_.SetDeadline(*c.rwdl)
    }
    if c.rdl != nil {
      c_.SetReadDeadline(*c.rdl)
    }
    if c.wdl != nil {
      c_.SetWriteDeadline(*c.wdl)
    }

    n, err := c_.Write(p)
    if err != nil {
      c.err = err
      return nil
    }

    if n != len(p) {
      c.err = fmt.Errorf("sendedn < len(p)")
      return nil
    }

    rbuf := make([]byte, 0, 1024)
    n, err = c_.Read(rbuf)
    if err != nil {
      c.err = err
      return nil
    }

    return rbuf[:n]
  }

  return nil
}

func (c *conn) resolveFromLocalDns(p []byte) {
  rData := c.query(p)
  if rData == nil {
    return
  }

  // in the case of etc
  if len(rData) == 1 {
    return
  }

  m := DNSMessageFromBytes(rData)
  if isHasErr(m.Header) {
    c.err = fmt.Errorf("header errCode: %d", m.Header.Flags&0b1111)
    return
  }
  if !isResponse(m.Header) {
    c.err = fmt.Errorf("not a response")
    return
  }
  if m.Header.ANCount != (uint16)(len(m.Answers)) {
    c.err = fmt.Errorf("ANcount: %d != len of answers: %d", m.Header.ANCount, len(m.Answers))
    return
  }
  c.rMsg.Answers = append(c.rMsg.Answers, m.Answers...)
  c.rMsg.AuthorityRRs = append(c.rMsg.AuthorityRRs, m.AdditionalRRs...)
  c.rMsg.AdditionalRRs = append(c.rMsg.AdditionalRRs, m.AdditionalRRs...)
}

func (c *conn) getRetValue(p []byte) {
  req := DNSMessageFromBytes(p)
  if len(req.Questions) == 0 {
    return
  }

  c.rMsg = req
  c.rMsg.Questions = req.Questions[:1]
  c.rMsg.Answers = req.Answers[:0]
  c.rMsg.AuthorityRRs = req.AuthorityRRs[:0]
  c.rMsg.AdditionalRRs = req.AdditionalRRs[:0]

  c.resolveFromLocalDns(p)
}

// Dial ...
func Dial(ctx context.Context, network, addr string) (net.Conn, error) {
  return newConn(ctx, network, addr)
}

// newConn ...
func newConn(_ context.Context, network, addr string) (net.Conn, error) {
  if _isSelfed {
    return &conn{
      network: network,
      addr:    addr,
    }, nil
  }

  var d net.Dialer
  return d.Dial(network, addr)
}

func Init(servers, domain string) {
  _nameserver = strings.Split(servers, ",")
  if domain != "" {
    _domain = domain
  }
  Resover.Dial = Dial

  http.DefaultTransport.(*http.Transport).DialContext = DialContent
}

var (
  _nameserver = []string{}
  _domain     string
  _isSelfed   bool = true

  Resover = net.Resolver{
    PreferGo: true,
    Dial:     nil,
  }
  Dialer = net.Dialer{
    Resolver: &Resover,
  }
  DialContent = func(ctx context.Context, network, addr string) (net.Conn, error) {
    return Dialer.DialContext(ctx, network, addr)
  }
  // conns      map[string]net.Conn // key = ip,
)
