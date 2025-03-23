package main

import (
	aterm_go "aterm-go"
	"fmt"
	"github.com/palantir/goastwriter/astgen"
	"github.com/palantir/goastwriter/decl"
	"github.com/palantir/goastwriter/expression"
	"github.com/palantir/goastwriter/statement"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	err := mainAux()
	if err != nil {
		log.Fatal(err)
	}
}

func mainAux() error {
	fset := new(token.FileSet)
	//f, _ := parser.ParseFile(fset, os.Args[1], nil, 0)
	ff := &decl.Function{
		Name: "Bar",
		FuncType: expression.FuncType{
			Params: []*expression.FuncParam{
				expression.NewFuncParam("input", expression.Type("Foo").Pointer()),
			},
			ReturnTypes: []expression.Type{
				expression.Type("Foo").Pointer(),
				expression.ErrorType,
			},
		},
		Body: []astgen.ASTStmt{
			&statement.Expression{
				Expr: expression.NewCallFunction("fmt", "Println"),
			},
			&statement.Expression{
				Expr: expression.NewCallFunction("gofmt", "Source", expression.Nil),
			},
			&statement.Return{
				Values: []astgen.ASTExpr{
					expression.VariableVal("input"),
					expression.Nil,
				},
			},
		},
	}
	_ = ff
	// parse the go file in /cmd/testdata/sum.go into f
	// here src: nil and mode: 0 means that the parser will read the file from disk
	// and will parse the whole file
	f, err := parser.ParseFile(fset, "cmd/testdata/sum.go", nil, 0)
	if err != nil {
		return err
	}
	//data, err := xml.Marshal(f)
	//if err != nil {
	//	panic(err)
	//}
	//_ = data
	//var v ast.FuncDecl
	//xml.Unmarshal(data, &v)
	//println(string(data))
	_ = fset
	x := ast.Ident{
		NamePos: 0,
		Name:    "Foo",
	}
	_ = x

	//ast.Fprint(os.Stdout, fset, x, aterm_go.NotPosInfoFilter)
	// create file /tmp/out.txt
	//fOut, err := os.Create("/tmp/out.txt")
	//ast.Fprint(fOut, fset, f, NotPosInfoFilter)
	_ = f
	//x := &ast.FieldList{
	//	List: []*ast.Field{
	//		&ast.Field{
	//			Names: []*ast.Ident{
	//				&ast.Ident{
	//					Name: "input",
	//				},
	//			},
	//			Type: &ast.Ident{
	//				Name: "&Foo",
	//			},
	//		},
	//	},
	//}
	marshal, err := aterm_go.MarshalWithFilter(f, aterm_go.NotPosInfoFilter)
	if err != nil {
		return err
	}
	fmt.Println(string(marshal))
	//_ = marshal
	//ast.Fprint(os.Stdout, fset, f, aterm_go.NotPosInfoFilter)
	//ast.Print(fset, f)
	return nil
}
