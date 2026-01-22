package svr

import (
  "go/token"
  "reflect"
  "strings"

  "github.com/yolksys/emei/log"
)

// ...
func parseRpc(rc_ any) error {
  tpy := reflect.TypeOf(rc_)
  tvalue := reflect.ValueOf(rc_)
  tname := reflect.Indirect(tvalue).Type().Name()
  if tname == "" {
    panic("tyepe name is empty")
  }

  if !token.IsExported(tname) {
    panic("is not exported type")
  }

  _rcvr := rcvrTx{
    params: make(map[string][]reflect.Type),
    funcs:  make(map[string]reflect.Value),
    value:  tvalue,
  }

  for i := range tpy.NumMethod() {
    m := tpy.Method(i)
    if !m.IsExported() {
      continue
    }

    mt_ := m.Type
    // at least have three params: receiver, *env and another
    if mt_.NumIn() < 2 {
      continue
    }

    if mt_.In(1) != typeOfEnv {
      continue
    }

    if mt_.NumIn() == 2 && (mt_.NumOut() == 0 ||
      (mt_.NumOut() == 1 && mt_.Out(0) == typeOfError)) {
      continue
    }

    isLegalF := false
    params := []reflect.Type{}
    for j := 2; j < mt_.NumIn(); j++ {
      pt_ := mt_.In(j)
      ptk := pt_.Kind()

      if ptk == reflect.Interface {
        if (j != mt_.NumIn()-1) ||
          ((pt_ != typeOfReader) &&
            (pt_ != typeOfWriter) &&
            (pt_ != typeOfReaderWriter)) {
          break
        }
      } else if ptk == reflect.Slice ||
        ptk == reflect.Array {
        pt_ = pt_.Elem()
        if pt_.Kind() == reflect.Pointer {
          pt_ = pt_.Elem()
        }
        if !isExportedStructOrBuiltinAtomType(pt_) {
          break
        }
      } else if ptk == reflect.Map {
        kpt := pt_.Key()
        pt_ = pt_.Elem()
        if pt_.Kind() == reflect.Pointer {
          pt_ = pt_.Elem()
        }
        if !isBuildinAtom(kpt) || !isExportedStructOrBuiltinAtomType(pt_) {
          break
        }
      } else if !isExportedStructOrBuiltinAtomType(pt_) {
        break
      }

      isLegalF = true
      params = append(params, pt_)
    }
    if !isLegalF {
      continue
    }

    mname := strings.ToLower(m.Name)
    _rcvr.funcs[mname] = m.Func
    _rcvr.params[mname] = params
  }

  if len(_rcvr.funcs) == 0 {
    panic("have no exported method")
  }

  log.Debug("*****", _rcvr)
  _rpcRecvs[strings.ToLower(tname)] = &_rcvr

  return nil
}

// Is this type exported or a builtin?
func isExportedStructOrBuiltinAtomType(t reflect.Type) bool {
  // PkgPath will be non-empty even for an exported type,
  // so we need to check the type name as well.
  return token.IsExported(t.Name()) ||
    isBuildinAtom(t)
}

// isBuildinAtom ...
func isBuildinAtom(t reflect.Type) bool {
  tk_ := t.Kind()

  return tk_ == reflect.Int ||
    tk_ == reflect.Int8 ||
    tk_ == reflect.Int16 ||
    tk_ == reflect.Int32 ||
    tk_ == reflect.Int64 ||
    tk_ == reflect.Uint ||
    tk_ == reflect.Uint8 ||
    tk_ == reflect.Uint16 ||
    tk_ == reflect.Uint32 ||
    tk_ == reflect.Uint64 ||
    tk_ == reflect.Float32 ||
    tk_ == reflect.Float64 ||
    tk_ == reflect.Complex64 ||
    tk_ == reflect.Complex128 ||
    tk_ == reflect.String
}

// getEeleTye ...
func getElemTye(typ reflect.Type) reflect.Type {
  // typ := reflect.TypeOf(v)
  k := typ.Kind()
  for {
    if k == reflect.Array ||
      k == reflect.Pointer ||
      k == reflect.Slice {
      typ = typ.Elem()
      k = typ.Kind()
      continue
    }
    return typ
  }
}

// getElemKind ...
func getElemKind(typ reflect.Type) reflect.Kind {
  // typ := reflect.TypeOf(v)
  k := typ.Kind()
  for {
    if k == reflect.Array ||
      k == reflect.Pointer ||
      k == reflect.Slice {
      typ = typ.Elem()
      k = typ.Kind()
      continue
    }
    return k
  }
}
