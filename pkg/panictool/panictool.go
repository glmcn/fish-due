package panictool

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

func RecoverWithError(err *error) {
	r := recover()
	if r != nil {
		if x, ok := r.(error); ok {
			*err = x
		} else {
			*err = errors.New(fmt.Sprint(r))
		}
	}
}

func String2Time(str string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	switch len(str) {
	case 19: // 2006-01-02 15:04:05
		if err != nil {
			return time.Parse(time.DateTime, str)
		}
		return time.ParseInLocation(time.DateTime, str, loc)
	case 10: // 2006-01-02
		if err != nil {
			return time.Parse(time.DateOnly, str)
		}
		return time.ParseInLocation(time.DateOnly, str, loc)
	case 8: // 15:04:05
		if err != nil {
			return time.Parse(time.TimeOnly, str)
		}
		return time.ParseInLocation(time.TimeOnly, str, loc)
	default:
		panic("invalid time string: " + str)
	}
}

func GetTimePointer[T string | *string | *time.Time](in T) *time.Time {
	t := any(in)
	switch v := t.(type) {
	case string:
		if v == "" {
			return nil
		}
		r, e := String2Time(v)
		if e != nil {
			panic(e)
		}
		return &r
	case *string:
		if v == nil {
			return nil
		}
		r, e := String2Time(*v)
		if e != nil {
			panic(e)
		}
		return &r
	case *time.Time:
		return v
	}
	panic(fmt.Errorf("invalid type [%s] for GetTimePointer", reflect.TypeOf(t)))
}

func GetString[T string | int8 | int | int32 | int64 | time.Time | *time.Time | *string | *int64](in T) string {
	t := any(in)
	switch v := t.(type) {
	case string:
		return v
	case int8:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case time.Time:
		return time.Time.Format(v, time.DateTime)
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return ""
	}

	switch v := t.(type) {
	case *string:
		return *v
	case *time.Time:
		return time.Time.Format(*v, time.DateTime)
	default:
		panic(fmt.Errorf("invalid type [%s] for GetString", reflect.TypeOf(t)))
	}
}

func GetStringPointer[T *string | *int8 | *int32 | *int64 | *time.Time | string | time.Time](in T) *string {
	t := any(in)
	switch v := (t).(type) {
	case string:
		return aws.String(v)
	case time.Time:
		return aws.String(GetString(v))
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return nil
	}

	switch v := (t).(type) {
	case *string:
		return aws.String(*v)
	case *int8:
		return aws.String(GetString(*v))
	case *int32:
		return aws.String(GetString(*v))
	case *int64:
		return aws.String(GetString(*v))
	case *time.Time:
		return aws.String(GetString(v))
	default:
		panic(fmt.Errorf("invalid type [%s] for GetStringPointer", reflect.TypeOf(t)))
	}
}

func GetInt64[T string | int8 | int32 | int64 | *string | *int8 | *int32 | *int64](in T) int64 {
	t := any(in)
	switch v := t.(type) {
	case string:
		if v == "" {
			return 0
		}
		r, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			panic(e)
		}
		return r
	case int8:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return 0
	}

	switch v := t.(type) {
	case *string:
		return GetInt64(*v)
	case *int8:
		return GetInt64(*v)
	case *int32:
		return GetInt64(*v)
	case *int64:
		return *v
	default:
		panic(fmt.Errorf("invalid type [%s] for GetInt64", reflect.TypeOf(t)))
	}
}

func GetInt32[T string | int8 | int32 | int64 | float64 | int | *string | *int8 | *int32 | *int64](in T) int32 {
	t := any(in)
	switch v := t.(type) {
	case string:
		if v == "" {
			return 0
		}
		r, e := strconv.Atoi(v)
		if e != nil {
			panic(e)
		}
		return int32(r)
	case int8:
		return int32(v)
	case int:
		return int32(v)
	case int32:
		return v
	case int64:
		if v != int64(int32(v)) {
			panic(fmt.Errorf("int64 %v convert to int32 failed for GetInt32", reflect.TypeOf(t)))
		}
		return int32(v)
	case float64:
		if v != float64(int32(v)) {
			panic(fmt.Errorf("float64 %v convert to int32 failed for GetInt32", reflect.TypeOf(t)))
		}
		return int32(v)
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return 0
	}

	switch v := t.(type) {
	case *string:
		return GetInt32(*v)
	case *int8:
		return GetInt32(*v)
	case *int32:
		return *v
	case *int64:
		return GetInt32(*v)
	default:
		panic(fmt.Errorf("invalid type [%s] for GetInt32", reflect.TypeOf(t)))
	}
}

func GetInt8[T string | int8 | int32 | int64 | *string | *int8 | *int32 | *int64](in T) int8 {
	t := any(in)
	switch v := t.(type) {
	case string:
		if v == "" {
			return 0
		}
		r, e := strconv.Atoi(v)
		if e != nil {
			panic(e)
		}
		return int8(r)
	case int8:
		return v
	case int32:
		if v != int32(int8(v)) {
			panic(fmt.Errorf("int32 %v convert to int8 failed for GetInt8", reflect.TypeOf(t)))
		}
		return int8(v)
	case int64:
		if v != int64(int8(v)) {
			panic(fmt.Errorf("int64 %v convert to int8 failed for GetInt8", reflect.TypeOf(t)))
		}
		return int8(v)
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return 0
	}

	switch v := t.(type) {
	case *string:
		return GetInt8(*v)
	case *int8:
		return *v
	case *int32:
		return GetInt8(*v)
	case *int64:
		return GetInt8(*v)
	default:
		panic(fmt.Errorf("invalid type [%s] for GetInt8", reflect.TypeOf(t)))
	}
}

func GetInt[T string | int8 | int | int32 | int64 | *string | *int8 | *int32 | *int | *int64](in T) int {
	t := any(in)
	switch v := t.(type) {
	case string:
		if v == "" {
			return 0
		}
		r, e := strconv.Atoi(v)
		if e != nil {
			panic(e)
		}
		return r
	case int8:
		return int(v)
	case int:
		return v
	case int32:
		if v != int32(int(v)) {
			panic(fmt.Errorf("int32 %v convert to int failed for GetInt", reflect.TypeOf(t)))
		}
		return int(v)
	case int64:
		if v != int64(int(v)) {
			panic(fmt.Errorf("int64 %v convert to int failed for GetInt", reflect.TypeOf(t)))
		}
		return int(v)
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return 0
	}

	switch v := t.(type) {
	case *string:
		return GetInt(*v)
	case *int8:
		return GetInt(*v)
	case *int32:
		return GetInt(*v)
	case *int64:
		return GetInt(*v)
	case *int:
		return GetInt(*v)
	default:
		panic(fmt.Errorf("invalid type [%s] for GetInt", reflect.TypeOf(t)))
	}
}

func GetInt64Pointer[T string | int8 | int32 | int64 | *string | *int8 | *int32 | *int64](in T) *int64 {
	t := any(in)
	switch v := t.(type) {
	case string:
		return aws.Int64(GetInt64(v))
	case int8:
		return aws.Int64(GetInt64(v))
	case int32:
		return aws.Int64(GetInt64(v))
	case int64:
		return aws.Int64(GetInt64(v))
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return nil
	}

	switch v := t.(type) {
	case *string:
		return aws.Int64(GetInt64(*v))
	case *int8:
		return aws.Int64(GetInt64(*v))
	case *int32:
		return aws.Int64(GetInt64(*v))
	case *int64:
		return aws.Int64(GetInt64(*v))
	default:
		panic(fmt.Errorf("invalid type [%s] for GetInt64Pointer", reflect.TypeOf(t)))
	}
}

func GetInt32Pointer[T string | int | int8 | int32 | int64 | *string | *int | *int8 | *int32 | *int64 | float64](in T) *int32 {
	t := any(in)
	switch v := t.(type) {
	case string:
		return aws.Int32(GetInt32(v))
	case int:
		return aws.Int32(GetInt32(v))
	case int8:
		return aws.Int32(GetInt32(v))
	case int32:
		return aws.Int32(GetInt32(v))
	case int64:
		return aws.Int32(GetInt32(v))
	case float64:
		return aws.Int32(GetInt32(v))
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return nil
	}

	switch v := t.(type) {
	case *string:
		return aws.Int32(GetInt32(*v))
	case *int8:
		return aws.Int32(GetInt32(*v))
	case *int:
		return aws.Int32(GetInt32(*v))
	case *int32:
		return aws.Int32(GetInt32(*v))
	case *int64:
		return aws.Int32(GetInt32(*v))
	default:
		panic(fmt.Errorf("invalid type [%s] for GetInt32Pointer", reflect.TypeOf(t)))
	}
}

func GetInt8Pointer[T string | int8 | int32 | int64 | *string | *int8 | *int32 | *int64](in T) *int8 {
	t := any(in)
	switch v := t.(type) {
	case string:
		return aws.Int8(GetInt8(v))
	case int8:
		return aws.Int8(GetInt8(v))
	case int32:
		return aws.Int8(GetInt8(v))
	case int64:
		return aws.Int8(GetInt8(v))
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return nil
	}

	switch v := t.(type) {
	case *string:
		return aws.Int8(GetInt8(*v))
	case *int8:
		return aws.Int8(GetInt8(*v))
	case *int32:
		return aws.Int8(GetInt8(*v))
	case *int64:
		return aws.Int8(GetInt8(*v))
	default:
		panic(fmt.Errorf("invalid type [%s] for GetInt8Pointer", reflect.TypeOf(t)))
	}
}

func GetIntPointer[T string | int8 | int32 | int64 | *string | *int8 | *int32 | *int64](in T) *int {
	t := any(in)
	switch v := t.(type) {
	case string:
		return aws.Int(GetInt(v))
	case int8:
		return aws.Int(GetInt(v))
	case int32:
		return aws.Int(GetInt(v))
	case int64:
		return aws.Int(GetInt(v))
	}

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return nil
	}

	switch v := t.(type) {
	case *string:
		return aws.Int(GetInt(*v))
	case *int8:
		return aws.Int(GetInt(*v))
	case *int32:
		return aws.Int(GetInt(*v))
	case *int64:
		return aws.Int(GetInt(*v))
	default:
		panic(fmt.Errorf("invalid type [%s] for GetIntPointer", reflect.TypeOf(t)))
	}
}

func GetBool[T *bool](in T) bool {
	t := any(in)

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return false
	}

	switch v := t.(type) {
	case *bool:
		return *v
	default:
		panic(fmt.Errorf("invalid type [%s] for GetBool", reflect.TypeOf(t)))
	}
}

func GetFloat64[T *int64](in T) float64 {
	t := any(in)

	// 上面过滤完成后剩下的都是指针
	if reflect.ValueOf(t).IsNil() {
		return 0
	}

	switch v := t.(type) {
	case *int64:
		return float64(*v)
	default:
		panic(fmt.Errorf("invalid type [%s] for GetFloat64", reflect.TypeOf(t)))
	}
}
