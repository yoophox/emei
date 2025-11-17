package env

func (m *Tjatse) TraceId() string {
  return m.Mid
}

func (m *Tjatse) TraceSpanId() string {
  return m.Sid
}

func (m *Tjatse) SetSID(i string) {
  m.Sid = i
}

func (m *Tjatse) SetTID(i string) {
  m.Mid = i
}

// User ...
func (m *Tjatse) Uid() string {
  return "" // jwd
}

func (m *Tjatse) UName() string {
  return "" // jwd
}
