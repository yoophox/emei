package env

const (
  DirResponse = "response"
  DirRequest  = "request"
  DirNotify   = "notify"
)

type Tjatse struct {
  // Len    uint64 `json:"len,omitempty"`
  Mid string `json:"mid,omitempty"` // 消息id, traced
  Jwt string `json:"jwt,omitempty"`
  Sid string `json:"sid,omitempty"` // span id
  // Met    string `json:"met,omitempty"` // call method
  Ver    string `json:"version,omitempty"`
  Typ    string `json:"typ,omitempty"`    // request， response，notify `request:"CustomRoleMgr"`...
  Code   int32  `json:"code,omitempty"`   // for response
  Reason string `json:"reason,omitempty"` // for response
}

// func Unmarshal(m []byte) (*Tjatse, error) {
// 	var r Tjatse
// 	e := json.Unmarshal(m, &r)

// 	if e != nil {
// 		//logger.Error("参数错误:", string(m))
// 		return nil, e
// 	}

// 	return &r
// }

// func (r *Tjatse) Write(w io.Writer, data ...interface{}) error {

// 	if r.Body == "" && len(data) > 0 {
// 		var _d interface{}
// 		if len(data) == 1 {
// 			_d = data[0]
// 		} else {
// 			_d = data
// 		}

// 		b, err := json.Marshal(_d)
// 		if err != nil {
// 			return errors.New("Fail: parm")
// 		}
// 		r.Body = b
// 	}

// 	b, err := json.Marshal(r)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = w.Write(b)
// 	return err
// }

// func (r *Tjatse) Check() error {
// 	return nil
// }

// func (r *Tjatse) GetUid() string {
// 	return ""
// }
