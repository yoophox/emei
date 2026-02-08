package ckube

import (
  "fmt"
  "reflect"

  "github.com/yoophox/emei/errs"
)

// assign ...
func assign(src, dst reflect.Value) (err error) {
  defer func() {
    r := recover()
    if r == nil {
      return
    }

    err = errs.ErrorF("err.cfg.encoder.struct.assign.panic", "%v", r)
  }()

  dst = getElemOfPointer(dst)
  f, err := getAssignFx(dst.Type())
  if err != nil {
    return err
  }

  return f(src, dst)
}

// assignArr ...
func assignArr(src, dst reflect.Value) error {
  if src.Kind() == reflect.Pointer {
    src = src.Elem()
  }

  if src.Kind() != reflect.Array {
    return errs.ErrorF("err.cfg.encoder.struct.assignArr", "source is not array")
  }
  _ = dst
  return nil
}

// assignMap ...
func assignMap(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignSlice(src, dst reflect.Value) error {
  if getElemOfPointer(src).Kind() != reflect.Slice {
    return errs.ErrorF("err.cfg.encoder.struct.assignSlice", "src is not slice")
  }

  if dst.IsZero() {
    dst = reflect.MakeSlice(dst.Elem().Type(), 0, src.Len())
  }
  eleK := dst.Type().Elem().Kind()
  var eleTyp reflect.Type = dst.Type().Elem()
  if eleK == reflect.Pointer {
    eleTyp = eleTyp.Elem()
  }

  var err error
  f, err := getAssignFx(eleTyp)
  for i := 0; i < src.Len(); i++ {
    ele := reflect.New(eleTyp)
    err = f(src.Index(i), ele)
    if err != nil {
      return err
    }

    if eleK != reflect.Pointer {
      ele = reflect.Append(dst, ele.Elem())
    } else {
      ele = reflect.Append(dst, ele)
    }
  }

  return nil
}

func assignStuct(src, dst reflect.Value) error {
  if getElemOfPointer(src).Kind() != reflect.Slice {
    return errs.ErrorF("err.cfg.encoder.struct.assignSlice", "src is not struct")
  }

  for i := 0; i < dst.NumField(); i++ {
    dstTyp := dst.Type()
    fSrc := src.FieldByName(dstTyp.Field(i).Name)
    if !fSrc.IsValid() {
      continue
    }

    var pfv reflect.Value
    fTyp := dst.Field(i).Type()
    fk_ := fTyp.Kind()
    if fk_ == reflect.Pointer {
      fTyp = fTyp.Elem()
      pfv = reflect.New(fTyp)
      dst.Field(i).Set(pfv)
    } else {
      pfv = dst.Field(i).Addr()
    }
    f, err := getAssignFx(fTyp)
    if err != nil {
      return err
    }
    err = f(fSrc, pfv.Elem())
    if err != nil {
      return err
    }
  }

  return nil
}

func assignBool(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignInt(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignInt8(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignInt16(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignInt32(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignInt64(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignFload32(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignFload64(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignComplex64(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignComplex128(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignUint(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignUint8(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignUint16(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignUint32(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

func assignUint64(src, dst reflect.Value) error {
  _, _ = src, dst
  return nil
}

// assignStr ...
func assignStr(src, dst reflect.Value) error {
  dst.SetString(fmt.Sprintf("%v", src.Interface()))
  return nil
}

// GetAssignFx ...
func getAssignFx(v reflect.Type) (assignFx, error) {
  switch v.Kind() {
  case reflect.Array:
    return assignArr, nil
  case reflect.Map:
    return assignMap, nil
  case reflect.Slice:
    return assignSlice, nil
  case reflect.Struct:
    return assignStuct, nil
  case reflect.Bool:
    return assignBool, nil
  case reflect.Int:
    return assignInt, nil
  case reflect.Int8:
    return assignInt8, nil
  case reflect.Int16:
    return assignInt16, nil
  case reflect.Int32:
    return assignInt32, nil
  case reflect.Int64:
    return assignInt64, nil
  case reflect.Float32:
    return assignFload32, nil
  case reflect.Float64:
    return assignFload64, nil
  case reflect.Complex64:
    return assignComplex64, nil
  case reflect.Complex128:
    return assignComplex128, nil
  case reflect.Uint:
    return assignUint, nil
  case reflect.Uint8:
    return assignUint8, nil
  case reflect.Uint16:
    return assignUint16, nil
  case reflect.Uint32:
    return assignUint32, nil
  case reflect.Uint64:
    return assignUint64, nil
  case reflect.String:
    return assignStr, nil
  }

  return nil, errs.ErrorF("err.cfg.encodet.struct.getAssignFx", "not supported type:%v", v)
}

type assignFx func(src, dst reflect.Value) error
