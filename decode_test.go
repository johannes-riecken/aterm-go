package aterm_go

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{name: "string", args: args{data: []byte(`"foo"`)}, want: "foo"},
		{name: "int", args: args{data: []byte(`42`)}, want: 42},
		{name: "[]int", args: args{data: []byte(`[1,2,3]`)}, want: []int{1, 2, 3}},
		{name: "Point", args: args{data: []byte(`Point(1,2)`)}, want: Point{1, 2}},
		{name: "Point3D", args: args{data: []byte(`Point3D(1,2,3)`)}, want: Point3D{to.Ptr(1), to.Ptr(2), to.Ptr(3)}},
	}

	var out0 string
	var out1 int
	var out2 []int
	var out3 Point
	var out4 Point3D

	t.Run(tests[0].name, func(t *testing.T) {
		if err := Unmarshal(tests[0].args.data, &out0); (err != nil) != tests[0].wantErr {
			t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tests[0].wantErr)
		}
		if !reflect.DeepEqual(out0, tests[0].want) {
			t.Errorf("Unmarshal() got = %v, want %v", out0, tests[0].want)
		}
	})

	t.Run(tests[1].name, func(t *testing.T) {
		if err := Unmarshal(tests[1].args.data, &out1); (err != nil) != tests[1].wantErr {
			t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tests[1].wantErr)
		}
		if !reflect.DeepEqual(out1, tests[1].want) {
			t.Errorf("Unmarshal() got = %v, want %v", out1, tests[1].want)
		}
	})

	t.Run(tests[2].name, func(t *testing.T) {
		if err := Unmarshal(tests[2].args.data, &out2); (err != nil) != tests[2].wantErr {
			t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tests[2].wantErr)
		}
		if !reflect.DeepEqual(out2, tests[2].want) {
			t.Errorf("Unmarshal() got = %v, want %v", out2, tests[2].want)
		}
	})

	t.Run(tests[3].name, func(t *testing.T) {
		if err := Unmarshal(tests[3].args.data, &out3); (err != nil) != tests[3].wantErr {
			t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tests[3].wantErr)
		}
		if !reflect.DeepEqual(out3, tests[3].want) {
			t.Errorf("Unmarshal() got = %v, want %v", out3, tests[3].want)
		}
	})

	t.Run(tests[4].name, func(t *testing.T) {
		if err := Unmarshal(tests[4].args.data, &out4); (err != nil) != tests[4].wantErr {
			t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tests[4].wantErr)
		}
		if !reflect.DeepEqual(out4, tests[4].want) {
			t.Errorf("Unmarshal() got = %v, want %v", out4, tests[4].want)
		}
	})

}
