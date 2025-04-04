package aterm_go

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
	"sort"
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

func encodeWithFilter(b *bytes.Buffer, v reflect.Value, filter func(_ string, v reflect.Value) bool) (err error, used bool) {
	if filter != nil && !filter("", v) {
		return nil, false
	}
	switch v.Kind() {
	case reflect.String:
		b.WriteString(strconv.Quote(v.Interface().(string)))
	case reflect.Int:
		if v.Type().String() != "int" {
			b.WriteString(fmt.Sprintf("%q", v.Interface()))
		} else {
			x := strconv.Itoa(int(v.Int()))
			b.WriteString(x)
		}
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
	case reflect.Map:
		err2, used2 := encodeMap(b, v, filter)
		if err2 != nil {
			return err2, used2
		}
	case reflect.Bool:
		// using golang.org/x/text/cases
		// strings.Title is replaced by the following code:
		// cases.Title(cases.Lower().String(v.Bool()))
		b.WriteString(cases.Title(language.English).String(strconv.FormatBool(v.Bool())))
	default:
		panic("unsupported type: " + v.Kind().String())
	}
	return nil, true
}

func encodeMap(b *bytes.Buffer, v reflect.Value, filter func(_ string, v reflect.Value) bool) (error, bool) {
	origLen := len(b.Bytes())
	b.WriteString(`"{`)
	keys := v.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	for i := 0; i < len(keys); i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(keys[i].String())
		b.WriteString(":=")
		err2, used2 := encodeWithFilter(b, v.MapIndex(keys[i]), filter)
		if err2 != nil {
			return err2, true
		}
		if !used2 {
			b.Truncate(origLen)
		}
	}
	if b.Len() > origLen {
		b.WriteString(`}"`)
	}
	return nil, false
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
	err, _ := encodeWithFilter(b, v, nil)
	return err
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
