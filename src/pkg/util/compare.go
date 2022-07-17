package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

const (
	COMPARE_GET_IMPLEMENT_VALUE_FUNCTION_NAME = "GetUtilCompareValue"
)

func IsNilInterfaceObject(i interface{}) bool {
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}

func compErrMsg(a, b interface{}) string {
	return fmt.Sprintf(
		"actual:\n%s\nexpect:\n%s\n---",
		printObj(a),
		printObj(b),
	)
}

func isNil(v *reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func compValueErrMsg(aV, bV *reflect.Value) string {
	var a, b interface{} = aV, bV
	if aV == nil || !aV.IsValid() {
		a = nil
	} else if aV.IsZero() {
		if isNil(aV) {
			a = nil
		} else {
			a = zeroValue(aV)
		}
	} else if aV.CanInterface() {
		a = aV.Interface()
	}

	if bV == nil || !bV.IsValid() {
		b = nil
	} else if bV.IsZero() {
		if isNil(bV) {
			b = nil
		} else {
			b = zeroValue(bV)
		}
	} else if bV.CanInterface() {
		b = bV.Interface()
	}
	return compErrMsg(a, b)
}

func zeroValue(v *reflect.Value) interface{} {
	var r interface{}
	k := v.Kind()
	switch k {
	case reflect.String:
		r = ""
	default:
		r = 0
	}
	return r
}

func printObj(a interface{}) string {
	b, _ := json.Marshal(a)
	return string(b)
}

func Comp(actual, expect interface{}) (bool, string) {
	if actual == nil {
		if expect == nil {
			return true, ""
		}
		return false, compErrMsg(actual, expect)
	}

	aV := reflect.ValueOf(actual)
	bV := reflect.ValueOf(expect)
	return compV(aV, bV)
}

func compV(aV, bV reflect.Value) (bool, string) {
	aK := aV.Kind()
	bK := bV.Kind()
	if aK != bK ||
		bV.IsZero() && !aV.IsZero() ||
		!bV.IsZero() && aV.IsZero() {
		return false, fmt.Sprintf(
			"actual:%s expect:%s\n%s",
			aK.String(),
			bK.String(),
			compValueErrMsg(&aV, &bV),
		)
	} else if bV.IsZero() && aV.IsZero() {
		return true, ""
	}

	switch bK {
	case reflect.Ptr:
		if aV.Pointer() == bV.Pointer() ||
			aV.IsNil() && bV.IsNil() {
			return true, ""
		}

		if !aV.IsNil() {
			aV = aV.Elem()
		}
		if !bV.IsNil() {
			bV = bV.Elem()
		}
		if !aV.CanInterface() || !bV.CanInterface() {
			return compV(aV, bV)
		}
		aVI := aV.Interface()
		bVI := bV.Interface()
		return Comp(aVI, bVI)
	case reflect.Slice, reflect.Array:
		if aK == reflect.Slice {
			if aV.IsNil() != bV.IsNil() {
				return false, compValueErrMsg(&aV, &bV)
			}
			if aV.Pointer() == bV.Pointer() {
				return true, ""
			}
		}

		aLen := aV.Len()
		bLen := bV.Len()
		maxLen := aLen
		if bLen > aLen {
			maxLen = bLen
		}

		errMsg := ""
		for i := 0; i < maxLen; i++ {
			var av, bv reflect.Value
			if i >= aLen {
				av = reflect.Value{}
				bv = bV.Index(i)
			} else if i >= bLen {
				av = aV.Index(i)
				bv = reflect.Value{}
			} else {
				av = aV.Index(i)
				bv = bV.Index(i)
			}
			if ok, msg := compV(av, bv); !ok {
				errMsg = fmt.Sprintf(
					"%s\n--- index: %d\n%s",
					errMsg,
					i,
					msg,
				)
			}
		}

		if errMsg != "" {
			return false, errMsg
		}
	case reflect.Struct:
		isGetImplementValue := false
		if actualImplementValue := aV.MethodByName(COMPARE_GET_IMPLEMENT_VALUE_FUNCTION_NAME); actualImplementValue.Kind() != reflect.Invalid {
			values := actualImplementValue.Call([]reflect.Value{})
			if len(values) == 1 {
				aV = values[0]
				isGetImplementValue = true
			}
		}
		if expectImplementValue := bV.MethodByName(COMPARE_GET_IMPLEMENT_VALUE_FUNCTION_NAME); expectImplementValue.Kind() != reflect.Invalid {
			values := expectImplementValue.Call([]reflect.Value{})
			if len(values) == 1 {
				bV = values[0]
				isGetImplementValue = true
			}
		}
		if isGetImplementValue {
			return compV(aV, bV)
		}

		if aV.NumField() != bV.NumField() {
			return false, compValueErrMsg(&aV, &bV)
		}

		var bt reflect.Type
		if !bV.CanInterface() {
			if aV.CanInterface() {
				return false, compValueErrMsg(&aV, nil)
			}
			bt = reflect.TypeOf(bV)
		} else if !aV.CanInterface() {
			return false, compValueErrMsg(nil, &bV)
		} else {
			actual := aV.Interface()
			expect := bV.Interface()
			if at, ok := actual.(time.Time); ok {
				if et, ok := expect.(time.Time); ok {
					if at.Equal(et) {
						return true, ""
					}
					return false, compErrMsg(at.String(), et.String())
				}
				return false, compErrMsg(at.String(), expect)
			}

			bt = reflect.TypeOf(expect)
		}

		for i := 0; i < bt.NumField(); i++ {
			fn := bt.Field(i).Name
			bv := bV.Field(i)
			av := aV.Field(i)

			if !bv.CanInterface() {
				if ok, msg := compV(av, bv); !ok {
					return false, fmt.Sprintf(
						"%s\n%s",
						fn,
						msg,
					)
				}
				continue
			}

			if ok, msg := Comp(av.Interface(), bv.Interface()); !ok {
				return false, fmt.Sprintf(
					"%s\n%s",
					fn,
					msg,
				)
			}
		}
	case reflect.Map:
		if aV.IsNil() != bV.IsNil() {
			return false, compValueErrMsg(&aV, &bV)
		}
		if aV.Pointer() == bV.Pointer() {
			return true, ""
		}

		if aV.Type().Key() != bV.Type().Key() ||
			aV.Type().Elem() != bV.Type().Elem() {
			return false, compValueErrMsg(&aV, &bV)
		}
		for _, ak := range aV.MapKeys() {
			bv := bV.MapIndex(ak)

			if !bv.IsValid() {
				av := aV.MapIndex(ak)
				return false, fmt.Sprintf(
					"%s\n%s",
					fmt.Sprint(ak.Interface()),
					compValueErrMsg(&av, nil),
				)
			}
		}
		for _, bk := range bV.MapKeys() {
			bv := bV.MapIndex(bk)
			av := aV.MapIndex(bk)

			if !bv.CanInterface() ||
				!av.IsValid() || !bv.IsValid() ||
				av.IsZero() && !bv.IsZero() || !av.IsZero() && bv.IsZero() {
				return false, fmt.Sprintf(
					"%s\n%s",
					fmt.Sprint(bk.Interface()),
					compValueErrMsg(&av, &bv),
				)
			}

			if ok, msg := Comp(av.Interface(), bv.Interface()); !av.IsValid() || !bv.IsValid() || !ok {
				return false, fmt.Sprintf(
					"%s\n%s",
					fmt.Sprint(bk.Interface()),
					msg,
				)
			}
		}
	default:
		var actual, expect interface{}
		switch bK {
		case
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			actual = aV.Int()
			expect = bV.Int()
		case
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			actual = aV.Uint()
			expect = bV.Uint()
		case
			reflect.Float32,
			reflect.Float64:
			actual = aV.Float()
			expect = bV.Float()
		case reflect.String:
			actual = aV.String()
			expect = bV.String()
		case reflect.Bool:
			actual = aV.Bool()
			expect = bV.Bool()
		case reflect.Interface:
			if aV.CanInterface() {
				actual = aV.Interface()
				expect = bV.Interface()
				return Comp(actual, expect)
			}
			return false, compValueErrMsg(nil, &bV)
		default:
			if !bV.CanInterface() {
				if aV.CanInterface() {
					return false, compValueErrMsg(&aV, nil)
				}
			}

			return false, ""
		}

		if !reflect.DeepEqual(actual, expect) {
			return false, compErrMsg(actual, expect)
		}
	}
	return true, ""
}

func CompareNilableToResult[Result any, T any](
	a, b T,
	compareFn func(a, b T) Result,
	aNilResult, bNilResult, bothNilResult Result,
) Result {
	var result Result
	CompareNilable(
		a, b,
		func() {
			result = aNilResult
		},
		func() {
			result = bNilResult
		},
		func() {
			result = compareFn(a, b)
		},
		func() {
			result = bothNilResult
		},
	)

	return result
}

func CompareNilableParseToResult[Result any, T any](
	a, b T,
	compareFn func(a, b T) Result,
	bothNilResult Result,
	getAResultFn func() Result,
	getBResultFn func() Result,
) Result {
	var result Result
	CompareNilable(
		a, b,
		func() {
			result = getBResultFn()
		},
		func() {
			result = getAResultFn()
		},
		func() {
			result = compareFn(a, b)
		},
		func() {
			result = bothNilResult
		},
	)

	return result
}

func CompareNilable[T any](
	a, b T,
	aNilFn func(),
	bNilFn func(),
	norNilFn func(),
	bothNilFn func(),
) {
	aIsZero := IsZero(a)
	bIsZero := IsZero(b)
	if aIsZero && bIsZero {
		bothNilFn()
	} else if aIsZero {
		aNilFn()
	} else if bIsZero {
		bNilFn()
	} else {
		norNilFn()
	}
}
