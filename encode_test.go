package aterm_go

import (
	"go/ast"
	"reflect"
	"testing"
)

type Point struct {
	X, Y int
}

type Nested struct {
	Point *Point
}

type OptPoint struct {
	X, Y *int
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
		{name: "Nested", args: args{x: Nested{Point: &Point{1, 2}}}, want: []byte(`Nested(Point(1,2))`)},
		{name: "OptPoint", args: args{x: OptPoint{X: new(int), Y: new(int)}}, want: []byte(`OptPoint(0,0)`)},
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

type Point3D struct {
	X, Y, Z *int
}

// copied from azcore
func Ptr[T any](v T) *T {
	return &v
}

func TestMarshalWithFilter(t *testing.T) {
	type args struct {
		x any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "Ident", args: args{x: ast.Ident{Name: "Foo"}}, want: []byte(`Ident("Foo")`)},
		{name: "struct with non-adjacent initialized fields", args: args{x: Point3D{X: Ptr[int](1), Z: Ptr[int](3)}}, want: []byte(`Point3D(1,3)`)},
		{name: "failing", args: args{x: ast.FieldList{List: []*ast.Field{
			{
				Names: []*ast.Ident{
					{
						Name: "foo",
					},
				},
			},
			{
				Names: []*ast.Ident{
					{
						Name: "bar",
					},
				},
			},
		}}}, want: []byte(`FieldList([Field([Ident("foo")]),Field([Ident("bar")])])`)},
		{name: "map[string]int", args: args{x: map[string]int{"foo": 1, "bar": 2}}, want: []byte(`"{bar:=2,foo:=1}"`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalWithFilter(tt.args.x, NotPosInfoFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalWithFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalWithFilter() got = %s, want %s", got, tt.want)
			}
		})
	}
}
