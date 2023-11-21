package aterm_go

import (
	"bytes"
	"reflect"
	"strconv"
)

type Foo struct {
}

func Marshal(x any) ([]byte, error) {
	v := reflect.ValueOf(x)
	var buf bytes.Buffer
	if err := encode(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(b *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		b.WriteString(strconv.Quote(v.Interface().(string)))
	case reflect.Int:
		x := strconv.Itoa(v.Interface().(int))
		b.WriteString(x)
	case reflect.Slice:
		b.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			if err := encode(b, v.Index(i)); err != nil {
				return err
			}
		}
		b.WriteByte(']')
	case reflect.Struct:
		typ := v.Type().String()[len(v.Type().PkgPath())+1:]
		b.WriteString(typ)
		b.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			encode(b, v.Field(i))
		}
		b.WriteByte(')')
	}
	return nil
}
