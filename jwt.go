package emei

import "github.com/yoophox/emei/jwt"

// type JWT interface {
// GetClaim(key string) string
// Exchange(opts ...Option) JWT
// IsLegal() bool
// Err() error
// Raw() string
// UID() string
// UNmae() string
// }
type JWT = jwt.JWT

var (
  // func New(o ...Option) JWT
  NewJwt = jwt.New
  // func FromStr(jwtstr string) JWT
  JwtFromStr = jwt.FromStr
  // func WithClaims(kv ...string) Option
  WithClaimsJwt = jwt.WithClaims
  // func WithExpireTime(etm time.Time) Option
  WithExpireJwt = jwt.WithExpireTime
  // func WithID(id string) Option
  JwtWithID = jwt.WithID
  // func WithIssuer(iss string) Option
  JwtWithIssuer = jwt.WithIssuer
  // func WithSubject(sub string) Option
  JwtWithSub = jwt.WithSubject
  // jwt.WithUserClaim(uid string, uname string)
  WithUserClaimJwt = jwt.WithUserClaim
)
