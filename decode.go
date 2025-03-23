package aterm_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

type lexer struct {
	data  []byte
	scan  scanner.Scanner
	token rune
}

func (l *lexer) next() {
	l.token = l.scan.Scan()
}

func (l *lexer) text() string {
	return l.scan.TokenText()
}

func (l *lexer) consume(want rune) error {
	if l.token != want {
		return fmt.Errorf("expected " + string(want) + ", got " + string(l.token))
	}
	l.next()
	return nil
}

func (l *lexer) endList() bool {
	return l.token == ']'
}

func (l *lexer) endStruct() bool {
	_ = json.Marshal
	return l.token == ')'
}

func UnmarshalWithSkips(data []byte, out interface{}, skips map[string][]int) error {
	lex := &lexer{data: data}
	lex.scan.Init(bytes.NewReader(data))
	lex.next()
	return readWithSkips(lex, reflect.ValueOf(out).Elem(), skips)
}

func readWithSkips(lex *lexer, v reflect.Value, skips map[string][]int) error {
	switch lex.token {
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		if v.Kind() == reflect.Ptr {
			v.Set(reflect.ValueOf(&i))
		} else {
			v.SetInt(int64(i))
		}
		lex.next()
	case '[':
		lex.next()
		err := readListWithSkips(lex, v, skips)
		if err != nil {
			return err
		}
		err = lex.consume(']')
		if err != nil {
			return err
		}
	case scanner.Ident: // struct name
		lex.next() // guess we don't need the struct name
		err := lex.consume('(')
		if err != nil {
			return err
		}
		err = readListWithSkips(lex, v, skips)
		if err != nil {
			return err
		}
		err = lex.consume(')')
		if err != nil {
			return err
		}
	default:
		panic("unhandled default case")
	}
	return nil
}

func readListWithSkips(lex *lexer, v reflect.Value, skips map[string][]int) error {
	switch v.Kind() {
	case reflect.Slice:
		err := readSliceWithSkips(lex, v, skips)
		if err != nil {
			return err
		}
	case reflect.Struct:
		err := readStructWithSkips(lex, v, skips)
		if err != nil {
			return err
		}
	default:
		panic("assertion error")
	}
	return nil
}

func readStructWithSkips(lex *lexer, v reflect.Value, skips map[string][]int) error {
	for i := 0; !lex.endStruct(); i++ {
		if i > 0 {
			err := lex.consume(',')
			if err != nil {
				return err
			}
		}
		currentSkips := skips[v.Type().Field(i).Name]
		if len(currentSkips) > 0 {
			for _, skip := range currentSkips {
				if skip == i {
					continue
				}
			}
		}
		err := readWithSkips(lex, v.FieldByIndex([]int{i}), skips)
		if err != nil {
			return err
		}
	}
	return nil
}

func readSliceWithSkips(lex *lexer, v reflect.Value, skips map[string][]int) error {
	for i := 0; !lex.endList(); i++ {
		if i > 0 {
			err := lex.consume(',')
			if err != nil {
				return err
			}
		}
		x := reflect.New(v.Type().Elem()).Elem()
		err := readWithSkips(lex, x, skips)
		if err != nil {
			return err
		}
		v.Set(reflect.Append(v, x))

	}
	return nil
}

func Unmarshal(data []byte, out interface{}) error {
	lex := &lexer{data: data}
	lex.scan.Init(bytes.NewReader(data))
	lex.next()
	return read(lex, reflect.ValueOf(out).Elem())
}

func read(lex *lexer, v reflect.Value) error {
	switch lex.token {
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		if v.Kind() == reflect.Ptr {
			v.Set(reflect.ValueOf(&i))
		} else {
			v.SetInt(int64(i))
		}
		lex.next()
	case '[':
		lex.next()
		err := readList(lex, v)
		if err != nil {
			return err
		}
		err = lex.consume(']')
		if err != nil {
			return err
		}
	case scanner.Ident: // struct name
		lex.next() // guess we don't need the struct name
		err := lex.consume('(')
		if err != nil {
			return err
		}
		err = readList(lex, v)
		if err != nil {
			return err
		}
		err = lex.consume(')')
		if err != nil {
			return err
		}
	default:
		panic("unhandled default case")
	}
	return nil
}

func readList(lex *lexer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Slice:
		err := readSlice(lex, v)
		if err != nil {
			return err
		}
	case reflect.Struct:
		err := readStruct(lex, v)
		if err != nil {
			return err
		}
	default:
		panic("assertion error")
	}
	return nil
}

func readStruct(lex *lexer, v reflect.Value) error {
	for i := 0; !lex.endStruct(); i++ {
		if i > 0 {
			err := lex.consume(',')
			if err != nil {
				return err
			}
		}
		err := read(lex, v.FieldByIndex([]int{i}))
		if err != nil {
			return err
		}
	}
	return nil
}

func readSlice(lex *lexer, v reflect.Value) error {
	for i := 0; !lex.endList(); i++ {
		if i > 0 {
			err := lex.consume(',')
			if err != nil {
				return err
			}
		}
		x := reflect.New(v.Type().Elem()).Elem()
		err := read(lex, x)
		if err != nil {
			return err
		}
		v.Set(reflect.Append(v, x))

	}
	return nil
}
