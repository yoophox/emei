package pki

import (
  "bytes"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha256"
  "crypto/tls"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/base64"
  "encoding/pem"
  "math/big"
  "os"
  "time"
)

// KeyPairWithPin returns PEM encoded Certificate and Key along with an SKPI
// fingerprint of the public key.
func KeyPairWithPin(bits int) (_cert []byte, _pri []byte, _pin []byte, err error) {
  privateKey, err := rsa.GenerateKey(rand.Reader, bits)
  if err != nil {
    return nil, nil, nil, err
  }

  tpl := x509.Certificate{
    SerialNumber:          big.NewInt(1),
    Subject:               pkix.Name{CommonName: "hzService_fwesadfwe"},
    NotBefore:             time.Now(),
    NotAfter:              time.Now().AddDate(2, 0, 0),
    BasicConstraintsValid: true,
    ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
    KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
  }
  derCert, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &privateKey.PublicKey, privateKey)
  if err != nil {
    return nil, nil, nil, err
  }

  buf := &bytes.Buffer{}
  err = pem.Encode(buf, &pem.Block{
    Type:  "CERTIFICATE",
    Bytes: derCert,
  })
  if err != nil {
    return nil, nil, nil, err
  }

  pemCert := buf.Bytes()

  buf = &bytes.Buffer{}
  err = pem.Encode(buf, &pem.Block{
    Type:  "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
  })
  if err != nil {
    return nil, nil, nil, err
  }
  pemKey := buf.Bytes()

  cert, err := x509.ParseCertificate(derCert)
  if err != nil {
    return nil, nil, nil, err
  }

  pubDER, err := x509.MarshalPKIXPublicKey(cert.PublicKey.(*rsa.PublicKey))
  if err != nil {
    return nil, nil, nil, err
  }
  sum := sha256.Sum256(pubDER)
  pin := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
  base64.StdEncoding.Encode(pin, sum[:])

  return pemCert, pemKey, pin, nil
}

func GenCertFile(dir string, bits ...int) error {
  _bits := 4096
  if len(bits) == 1 {
    _bits = bits[0]
  }

  cert, key, pin, err := KeyPairWithPin(_bits)
  if err != nil {
    return err
  }

  var fcert, fkey, fpin *os.File

  defer func() {
    if fcert != nil {
      fcert.Close()
    }
    if fkey != nil {
      fcert.Close()
    }
    if fpin != nil {
      fcert.Close()
    }
  }()

  fcert, err0 := os.OpenFile(dir+"/cert.pem", os.O_CREATE|os.O_WRONLY, 0777)
  fkey, err1 := os.OpenFile(dir+"/key.pem", os.O_CREATE|os.O_WRONLY, 0777)
  fpin, err4 := os.OpenFile(dir+"/cert.pin", os.O_CREATE|os.O_WRONLY, 0777)

  if err0 != nil ||
    err1 != nil ||
    err4 != nil {
    return err
  }

  _, err0 = fcert.Write(cert)
  _, err1 = fkey.Write(key)
  _, err4 = fpin.Write(pin)
  if err0 != nil ||
    err1 != nil ||
    err4 != nil {
    return err
  }

  return nil
}

func GenCert(bits ...int) (cert []byte, key []byte, pin []byte, err error) {
  _bits := 4096
  if len(bits) == 1 {
    _bits = bits[0]
  }

  cert, key, pin, err = KeyPairWithPin(_bits)
  return
}

// func GenTlsCfg(isSvr bool) (*tls.Config, error) {
//   var mtls bool
//   e := cfg.GetCfgItem("net.mtls", &mtls)
//   if (e != nil || !mtls) && !isSvr {
//     return GenCltTlsCfg(), nil
//   }

//   c, e := GenSvrTlsCfg()
//   if e != nil {
//     return nil, e
//   }
//   if e != nil || !mtls {
//     return c, nil
//   }

//   var caP string
//   e = cfg.GetCfgItem("net.crypto.ca", &caP)
//   if e != nil {
//     return c, nil
//   }
//   ca, e := os.ReadFile(caP)
//   if e != nil {
//     return c, nil
//   }
//   CAs := x509.NewCertPool()
//   ok := CAs.AppendCertsFromPEM(ca)
//   if !ok {
//     return c, nil
//   }
//   opt := x509.VerifyOptions{
//     DNSName:       "yolk.com",
//     Intermediates: x509.NewCertPool(),
//   }

//   vf := func(cs tls.ConnectionState) error {
//     for _, v := range cs.PeerCertificates[1:] {
//       opt.Intermediates.AddCert(v)
//     }

//     _, e := cs.PeerCertificates[0].Verify(opt)
//     return e
//   }

//   if isSvr {
//     c.ClientCAs = CAs
//     c.ClientAuth = tls.RequireAnyClientCert
//     opt.KeyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
//     c.VerifyConnection = vf
//   } else {
//     c.RootCAs = CAs
//     c.InsecureSkipVerify = true
//     c.VerifyConnection = vf
//   }

//   return c, nil
// }

// func GenSvrTlsCfg() (*tls.Config, error) {
//   var certP, keyP string
//   e := cfg.GetCfgItem("net.crypto.cert", &certP)
//   e0 := cfg.GetCfgItem("net.crypto.key", &keyP)
//   if e != nil || e0 != nil {
//     return nil, errors.New("cert or key error")
//   }

//   cert, e := os.ReadFile(certP)
//   if e != nil {
//     return nil, errors.New("read cert error")
//   }
//   key, e := os.ReadFile(keyP)
//   if e != nil {
//     return nil, errors.New("read key error")
//   }

//   tlsCert, e := tls.X509KeyPair(cert, key)
//   if e != nil {
//     return nil, e
//   }

//   c := &tls.Config{
//     Certificates: []tls.Certificate{tlsCert},
//     NextProtos:   []string{"asderfdfkwekjasdfwe"},
//   }

//   return c, nil
// }

func GenCltTlsCfg() *tls.Config {
  return &tls.Config{
    // Certificates: []tls.Certificate{tlsCert},
    NextProtos:         []string{"asderfdfkwekjasdfwe"},
    InsecureSkipVerify: true,
  }
}
