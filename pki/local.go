package pki

//@key:[any space]  {type:<key type(ed25519/rsa...)>}.{cid:<id>}.{pri-len:<byts num>}.{pub-len:}\n
//<prikey>\n<pubkey>\n
//can't have any other letter before @key

// pem format
// public key:pem of PKIX, ASN.1 DER
// private key: pem of PKCS #8, ASN.1 DER

import (
  "bufio"
  "crypto/ed25519"
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "io"
  "math/rand"
  "os"
  "regexp"
  "strconv"

  "github.com/fsnotify/fsnotify"
  "github.com/yoophox/emei/errs"
)

type local struct {
  cryptos map[uint64]*cryptosx
  crps    map[string][]*cryptosx // key = type: rsa
}

// newLocal ...
func newLocalBackend() *local {
  _, _ = fsnotify.NewWatcher()
  f, err := os.Open(_pkiLocal)
  if err != nil {
    panic(err)
  }
  defer f.Close()
  regS := regexp.MustCompile(`@key.*\{type:([0-9a-zA-Z]+)\}\.\{cid:([0-9]+)\}\.\{pri-len:([0-9]+)\}\.\{pub-len:([0-9]+)\}`)
  // regS := regexp.MustCompile(`^@key.*type:[a-z]+`)

  fsc := bufio.NewReader(f)
  var c *cryptosx
  var pril, publ int

  l := &local{cryptos: map[uint64]*cryptosx{}, crps: map[string][]*cryptosx{}}
  for line, err := fsc.ReadBytes('\n'); ; line, err = fsc.ReadBytes('\n') {
    if err != nil {
      if err == io.EOF {
        break
      }

      panic(err)
    }
    if len(line) == 0 {
      continue
    }

    values := regS.FindAllStringSubmatch(string(line), -1)
    if values == nil {
      panic(string(line))
    }
    ct_ := values[0][1]
    cid := values[0][2]
    ci_ := values[0][3]
    cu_ := values[0][4]
    if string(ct_) == "" || string(ci_) == "" || string(cu_) == "" {
      panic("format error")
    }

    c = &cryptosx{}
    c.Typ = string(ct_)
    cid_, err := strconv.ParseInt(string(cid), 10, 64)
    if err != nil {
      panic("parse cid" + string(cid))
    }
    c.CID = uint64(cid_)
    pril, err = strconv.Atoi(ci_)
    if err != nil {
      panic("parse pril")
    }
    publ, err = strconv.Atoi(cu_)
    if err != nil {
      panic("parse publ")
    }
    ksc := make([]byte, pril+1+publ+1)
    n, err := io.ReadFull(fsc, ksc)
    if n != (pril + 1 + publ + 1) {
      fmt.Println("************", n, pril, publ, len(ksc))
      panic(err)
    }
    // var prib, pubb pem.Block
    prib, _ := pem.Decode(ksc[:pril])
    pri, err := x509.ParsePKCS8PrivateKey(prib.Bytes)
    if err != nil {
      panic(err)
    }
    pubb, _ := pem.Decode(ksc[pril+1 : len(ksc)-1])
    pub, err := x509.ParsePKIXPublicKey(pubb.Bytes)
    if err != nil {
      panic(err)
    }
    c.Pri = pri
    c.Pub = pub
    _, ok := l.crps[c.Typ]
    if !ok {
      l.crps[c.Typ] = []*cryptosx{}
    }
    l.crps[c.Typ] = append(l.crps[c.Typ], c)
    l.cryptos[uint64(cid_)] = c

    c = nil
    pril = 0 //
    publ = 0
  }

  return l
}

func (l *local) getRandomCrypto(o ...Option) (uint64, error) {
  var opt options

  for _, v := range o {
    v(&opt)
  }

  typ := opt.typ
  typc, ok := l.crps[string(typ)]
  if !ok || len(typc) <= 0 {
    return 0, errs.Wrap(fmt.Errorf("have no keys for type:%s in local data", typ), ERR_ID_PKI_LOCAL_NO_KEY)
  }

  i := rand.Intn(len(typc))

  return typc[i].CID, nil
}

func (l *local) getPriKeyByID(jwt string, id uint64) (*pem.Block, error) {
  return nil, nil
}

func (l *local) sign(id uint64, c []byte) ([]byte, error) {
  cyt, ok := l.cryptos[id]
  if !ok {
    return nil, errs.Wrap(fmt.Errorf("have no key for id:%d", id), ERR_ID_PKI_LOCAL_NO_KEY)
  }

  switch cyt.Typ {
  case "rsa":
  case "ed25519":
    ret := ed25519.Sign(cyt.Pri.(ed25519.PrivateKey), c)
    return ret, nil
  }
  return nil, errs.Wrap(fmt.Errorf("unkown typ of key:%s", cyt.Typ), ERR_ID_PKI_LOCAL_TYP)
}

func (l *local) getPubKeyByID(id uint64) (*pem.Block, error) {
  return nil, nil
}
