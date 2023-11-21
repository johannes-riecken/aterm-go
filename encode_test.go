package aterm_go

import (
	"bytes"
	"reflect"
	"testing"
)

type Point struct {
	X, Y int
}

func TestMarshal(t *testing.T) {
	type args struct {
		x any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "string", args: args{x: "foo"}, want: []byte(`"foo"`)},
		{name: "int", args: args{x: 42}, want: []byte(`42`)},
		{name: "[]int", args: args{x: []int{1, 2, 3}}, want: []byte(`[1,2,3]`)},
		{name: "Point", args: args{x: Point{1, 2}}, want: []byte(`Point(1,2)`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_encode(t *testing.T) {
	type args struct {
		b *bytes.Buffer
		v reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := encode(tt.args.b, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
