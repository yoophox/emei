package resolver

import (
  "bytes"
)

// DNSMessage represents a DNS message.
// It contains the header, questions, answers, authority RRs, and additional RRs.
//
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1 for more information
type DNSMessage struct {
  Header        *Header
  Questions     []*Question
  Answers       []*RR
  AuthorityRRs  []*RR
  AdditionalRRs []*RR
}

// NewDNSMessage creates a new DNSMessage with the given header, questions, and resource records.
// It returns a pointer to the created DNSMessage.
func NewDNSMessage() *DNSMessage {
  return &DNSMessage{}
}

// ToBytes converts the DNSMessage to a byte slice.
// It returns the byte slice representation of the DNSMessage.
func (m *DNSMessage) ToBytes(buf_ ...[]byte) []byte {
  // Create a buffer to store the bytes
  var buf *bytes.Buffer

  if len(buf_) > 0 {
    buf = bytes.NewBuffer(buf_[0])
  } else {
    buf = new(bytes.Buffer)
  }

  // Write the header to the buffer
  buf.Write(m.Header.ToBytes())

  // Write the questions to the buffer
  for _, q := range m.Questions {
    buf.Write(q.ToBytes())
  }

  // Write the answers to the buffer
  for _, a := range m.Answers {
    buf.Write(a.ToBytes())
  }

  // Write the authority RRs to the buffer
  for _, rr := range m.AuthorityRRs {
    buf.Write(rr.ToBytes())
  }

  // Write the additional RRs to the buffer
  for _, rr := range m.AdditionalRRs {
    buf.Write(rr.ToBytes())
  }

  // Return the bytes from the buffer
  return buf.Bytes()
}

// appendFromBufferUntilNull reads bytes from the buffer until a null byte is encountered.
// It returns the read bytes as a byte slice.
func appendFromBufferUntilNull(buf *bytes.Buffer) []byte {
  // Create a bytes slice by reading the bytes until we reach a null byte for any string field
  data := make([]byte, 0)
  for {
    b := buf.Next(1)
    data = append(data, b[0])
    if b[0] == 0 {
      break
    }
  }
  return data
}

// DNSMessageFromBytes creates a DNSMessage from the given byte slice.
// It returns a pointer to the created DNSMessage.
func DNSMessageFromBytes(data []byte) *DNSMessage {
  // Create a new buffer from the data
  buf := bytes.NewBuffer(data)
  bufCopy := bytes.NewBuffer(data)

  // Read the header from the buffer
  header := HeaderFromBytes(buf.Next(12))

  // Read the questions from the buffer
  questions := make([]*Question, header.QDCount)
  for i := 0; i < int(header.QDCount); i++ {
    questionBytes := appendFromBufferUntilNull(buf)
    questionBytes = append(questionBytes, buf.Next(4)...)
    questions[i] = QuestionFromBytes(questionBytes)
  }

  // Read the answers from the buffer
  answers := make([]*RR, header.ANCount)
  for i := 0; i < int(header.ANCount); i++ {
    rrBytes := TrimRRBytes(buf)
    answers[i] = RRFromBytes(rrBytes, bufCopy)
  }

  // Read the authority RRs from the buffer
  authorityRRs := make([]*RR, header.NSCount)
  for i := 0; i < int(header.NSCount); i++ {
    rrBytes := TrimRRBytes(buf)
    authorityRRs[i] = RRFromBytes(rrBytes, bufCopy)
  }

  // Read the additional RRs from the buffer
  additionalRRs := make([]*RR, header.ARCount)
  for i := 0; i < int(header.ARCount); i++ {
    rrBytes := TrimRRBytes(buf)
    additionalRRs[i] = RRFromBytes(rrBytes, bufCopy)
  }

  // Create a new DNS message with the parsed data
  return &DNSMessage{
    Header:        header,
    Questions:     questions,
    Answers:       answers,
    AuthorityRRs:  authorityRRs,
    AdditionalRRs: additionalRRs,
  }
}
