package core

import (
  "fmt"
  "time"

  "github.com/yolksys/emei/log/cache"
  "github.com/yolksys/emei/log/fmt/json"
)

func serialize(buf []byte, cont ...interface{}) cache.Content {
  b := buf
  l := len(cont)
  l = l - l&1
  for i := 0; i < l; i += 2 {
    key, ok := cont[i].(string)
    if !ok {
      continue
    }

    switch m := cont[i+1].(type) {
    case string:
      b = enc.AppendString(enc.AppendKey(b, key), m)
      break
    case []string:
      b = enc.AppendStrings(enc.AppendKey(b, key), m)
      break
    case fmt.Stringer:
      b = enc.AppendStringer(enc.AppendKey(b, key), m)
      break
    case []fmt.Stringer:
      b = enc.AppendStringers(enc.AppendKey(b, key), m)
      break
    case []byte:
      b = enc.AppendBytes(enc.AppendKey(b, key), m)
      break
    case error:
      b = enc.AppendString(enc.AppendKey(b, key), m.Error())
      break
    case []error:
      enc.AppendKey(b, key)
      b = append(b, '[')
      for _, v := range m {
        b = enc.AppendString(append(b, ','), v.Error())
      }
      b = append(b, ']')

      break
    case bool:
      b = enc.AppendBool(enc.AppendKey(b, key), m)
      break
    case []bool:
      b = enc.AppendBools(enc.AppendKey(b, key), m)
      break
    case int:
      b = enc.AppendInt(enc.AppendKey(b, key), m)
      break
    case []int:
      b = enc.AppendInts(enc.AppendKey(b, key), m)
      break
    case int8:
      b = enc.AppendInt8(enc.AppendKey(b, key), m)
      break
    case []int8:
      b = enc.AppendInts8(enc.AppendKey(b, key), m)
      break
    case int16:
      b = enc.AppendInt16(enc.AppendKey(b, key), m)
      break
    case []int16:
      b = enc.AppendInts16(enc.AppendKey(b, key), m)
      break
    case int32:
      b = enc.AppendInt32(enc.AppendKey(b, key), m)
      break
    case []int32:
      b = enc.AppendInts32(enc.AppendKey(b, key), m)
      break
    case int64:
      b = enc.AppendInt64(enc.AppendKey(b, key), m)
      break
    case []int64:
      b = enc.AppendInts64(enc.AppendKey(b, key), m)
      break
    case uint:
      b = enc.AppendUint(enc.AppendKey(b, key), m)
      break
    case []uint:
      b = enc.AppendUints(enc.AppendKey(b, key), m)
      break
    case uint8:
      b = enc.AppendUint8(enc.AppendKey(b, key), m)
      break
    case uint16:
      b = enc.AppendUint16(enc.AppendKey(b, key), m)
      break
    case []uint16:
      b = enc.AppendUints16(enc.AppendKey(b, key), m)
      break
    case uint32:
      b = enc.AppendUint32(enc.AppendKey(b, key), m)
      break
    case []uint32:
      b = enc.AppendUints32(enc.AppendKey(b, key), m)
      break
    case uint64:
      b = enc.AppendUint64(enc.AppendKey(b, key), m)
      break
    case []uint64:
      b = enc.AppendUints64(enc.AppendKey(b, key), m)
      break
    case float32:
      b = enc.AppendFloat32(enc.AppendKey(b, key), m, FloatingPointPrecision)
      break
    case []float32:
      b = enc.AppendFloats32(enc.AppendKey(b, key), m, FloatingPointPrecision)
      break
    case float64:
      b = enc.AppendFloat64(enc.AppendKey(b, key), m, FloatingPointPrecision)
      break
    case []float64:
      b = enc.AppendFloats64(enc.AppendKey(b, key), m, FloatingPointPrecision)
      break
    case time.Duration:
      break
    default:
      b = enc.AppendInterface(enc.AppendKey(b, key), m)
    }
  }
  return b
}

var (
  // FloatingPointPrecision, if set to a value other than -1, controls the number
  // of digits when formatting float numbers in JSON. See strconv.FormatFloat for
  // more details.
  FloatingPointPrecision = -1
  enc                    = json.Encoder{}
)
