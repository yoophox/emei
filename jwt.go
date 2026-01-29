package emei

import "github.com/yoophox/emei/jwt"

var (
  JwtNew        = jwt.New
  JwtFromStr    = jwt.FromStr
  JwtWithClaims = jwt.WithClaims
  JwtWithExpire = jwt.WithExpireTime
  JwtWithID     = jwt.WithID
  JwtWithIssuer = jwt.WithIssuer
  JwtWithSub    = jwt.WithSubject
)
