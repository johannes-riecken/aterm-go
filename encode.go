package aterm_go

import (
	"bytes"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

var seen = make(map[any]bool)

func Marshal(x any) ([]byte, error) {
	v := reflect.ValueOf(x)
	var buf bytes.Buffer
	if err := encode(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalWithFilter(x any, filter func(_ string, v reflect.Value) bool) ([]byte, error) {
	v := reflect.ValueOf(x)
	var buf bytes.Buffer
	if err, _ := encodeWithFilter(&buf, v, filter); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Add() []byte {
	b := bytes.NewBuffer([]byte("foo"))
	b2 := bytes.NewBuffer(b.Bytes()[:len(b.Bytes())-1])
	b2.WriteByte('b')
	b2.WriteByte('a')
	b2.WriteByte('r')
	return b2.Bytes()
}

func encodeWithFilter(b *bytes.Buffer, v reflect.Value, filter func(_ string, v reflect.Value) bool) (err error, used bool) {
	if !filter("", v) {
		return nil, false
	}
	switch v.Kind() {
	case reflect.String:
		b.WriteString(strconv.Quote(v.Interface().(string)))
	case reflect.Int:
		x := strconv.Itoa(int(v.Int()))
		b.WriteString(x)
	case reflect.Slice:
		iEnd := v.Len()
		parenChars := []byte("[]")
		selector := func(v reflect.Value, i int) reflect.Value {
			return v.Index(i)
		}
		err, used := encodeCommaSeparated(b, v, filter, selector, parenChars, iEnd)
		if err != nil {
			return err, used
		}
	case reflect.Struct:
		err2, used2 := encodeStruct(b, v, filter)
		if err2 != nil {
			return err2, used2
		}
	case reflect.Pointer:
		if seen[v.Interface()] {
			return nil, false
		}
		seen[v.Interface()] = true
		return encodeWithFilter(b, v.Elem(), filter)
	case reflect.Invalid:
		b.WriteString("nil")
	case reflect.Interface:
		return encodeWithFilter(b, v.Elem(), filter)
	default:
		panic("unsupported type: " + v.Kind().String())
	}
	return nil, true
}

func encodeStruct(b *bytes.Buffer, v reflect.Value, filter func(_ string, v reflect.Value) bool) (error, bool) {
	pkgPath := v.Type().PkgPath()
	pkgPath = pkgPath[strings.LastIndex(pkgPath, "/")+1:]
	typ := v.Type().String()[len(pkgPath)+1:]
	b.WriteString(typ)
	iEnd := v.NumField()
	parenChars := []byte("()")
	selector := func(v reflect.Value, i int) reflect.Value {
		return v.Field(i)
	}
	return encodeCommaSeparated(b, v, filter, selector, parenChars, iEnd)
}

func encodeCommaSeparated(b *bytes.Buffer, v reflect.Value, filter func(_ string, v reflect.Value) bool, selector func(v reflect.Value, i int) reflect.Value, parenChars []byte, iEnd int) (error, bool) {
	b.WriteByte(parenChars[0])
	origLen := len(b.Bytes())
	needsComma := false
	for i := 0; i < iEnd; i++ {
		if needsComma {
			b.WriteByte(',')
		}
		err2, used2 := encodeWithFilter(b, selector(v, i), filter)
		if err2 != nil {
			return err2, true
		}
		if !used2 {
			if len(b.Bytes()) > origLen {
				b.Truncate(len(b.Bytes()) - 1)
			}
		} else {
			needsComma = true
		}
	}
	b.WriteByte(parenChars[1])
	return nil, false
}

func encode(b *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		b.WriteString(strconv.Quote(v.Interface().(string)))
	case reflect.Int:
		x := strconv.Itoa(int(v.Int()))
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
		// we can't use `typ := v.Type().String()[len(v.Type().PkgPath())+1:]`
		// because PkgPath() is the full path, not just the package name
		// better solution:
		pkgPath := v.Type().PkgPath()
		// only keep the last part of the path
		pkgPath = pkgPath[strings.LastIndex(pkgPath, "/")+1:]
		typ := v.Type().String()[len(pkgPath)+1:]
		b.WriteString(typ)
		b.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			encode(b, v.Field(i))
		}
		b.WriteByte(')')
	case reflect.Pointer:
		return encode(b, v.Elem())
	case reflect.Invalid:
		b.WriteString("nil")
	case reflect.Interface:
		return encode(b, v.Elem())
	default:
		panic("unsupported type: " + v.Kind().String())
	}
	return nil
}

var NotPosInfoFilter ast.FieldFilter = func(name string, value reflect.Value) bool {
	if !ast.NotNilFilter(name, value) {
		return false
	}
	if _, ok := value.Interface().(token.Pos); ok {
		return false
	}
	return true
}
