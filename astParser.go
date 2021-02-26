package KlaudiasGoTrace

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strconv"
	"strings"
)

func toString(fset *token.FileSet, file interface{}) string {
	var buffer bytes.Buffer
	printer.Fprint(&buffer, fset, file)
	return buffer.String()
}

func repr(str string) string {
	return fmt.Sprintf("%#v", str)
}

func printTree(fset *token.FileSet, file *ast.File) {
	//printer.Fprint(os.Stdout, fset, file)
	fmt.Println(toString(fset, file))
}

func newCallExpr(exp string) *ast.CallExpr {

	funex, _ := parser.ParseExpr(exp)
	expr := ast.CallExpr{
		Fun:      funex,
		Lparen:   0,
		Args:     nil,
		Ellipsis: 0,
		Rparen:   0,
	}
	return &expr
}

func newField(name string, typ string) *ast.Field {
	field := &ast.Field{
		Doc:     nil,
		Names:   []*ast.Ident{ast.NewIdent(name)},
		Type:    ast.NewIdent(typ),
		Tag:     nil,
		Comment: nil,
	}
	return field
}

func insert(a []ast.Stmt, index int, value ast.Stmt) []ast.Stmt {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func prepend(data []ast.Stmt, value ast.Stmt) []ast.Stmt {
	data = append([]ast.Stmt{value}, data...)
	return data
}

func getFuncName(line string) string {
	splitted := strings.Split(line, "(")
	if len(splitted) >= 2 {
		splitted2 := strings.Split(splitted[0], " ")
		if len(splitted) >= 2 {
			return splitted2[1]
		}
	}
	return ""
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func main() {
	// parsovanie suboru do AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test2.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var goFuncs []string

	// prehladavanie AST
	ast.Inspect(node, func(n ast.Node) bool {
		fungo, ok := n.(*ast.GoStmt)
		if ok {
			name := getFuncName(toString(fset, fungo))
			goFuncs = append(goFuncs, name)
		}
		return true
	})

	ast.Inspect(node, func(n ast.Node) bool {
		// pridanie argumentu parentId do vsetkych funkcii
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			field := newField("mnouDodanyParameter123", "uint64")
			fn.Type.Params.List = append(fn.Type.Params.List, field)
		}

		// hladanie vsetkych deklaraácii funkcii, ktore vystupuju ako gorutiny
		fn2, ok := n.(*ast.FuncDecl)
		if ok && contains(goFuncs, fn.Name.Name) {
			fmt.Println(strconv.Itoa(fset.Position(fn.Pos()).Line), fn2.Name.Name, len(fn2.Type.Params.List))
			//field := newField("mnouDodanyParameter123", "uint64")
			//fn.Type.Params.List = append(fn.Type.Params.List, field)

			////vlozenie vyrazu na zaciatok funkcie
			expr, _ := parser.ParseExpr("KlaudiasGoTrace.StartGoroutine(parentId)")
			stmt := ast.ExprStmt{X: expr}
			fn.Body.List = prepend(fn.Body.List, &stmt)

			////vlozenie vyrazu na koniec funkcie
			expr2, _ := parser.ParseExpr("KlaudiasGoTrace.StopGoroutine()")
			stmt2 := ast.ExprStmt{X: expr2}
			fn.Body.List = append(fn.Body.List, &stmt2)

			// vypisanie vsetkych parametrov funkcie
			//for i := 0; i < len(fn.Type.Params.List); i++ {
			//	fmt.Println(fn.Type.Params.List[i].Names[0].Name)
			//}
		}
		//hladanie vsetkych volani
		//callEx, ok := n.(*ast.CallExpr)
		//if ok {
		//	fmt.Println(strconv.Itoa(fset.Position(callEx.Pos()).Line), callEx.Fun, callEx.End())
		//
		//}

		////hladanie vsetkych go volani
		//fungo, ok := n.(*ast.GoStmt)
		//if ok {
		//	//fmt.Println(strconv.Itoa(fset.Position(fungo.Pos()).Line), fungo.Call.Fun)
		//	//str, _ := fmt.Println(fungo.Call.Fun)
		//	//fmt.Println(str)
		//	//goFuncs = append(goFuncs, fungo.Call.Fun)
		//	fmt.Println(reflect.TypeOf(fungo.Call.))
		//}

		//// hladanie vsetkych posielani kanalom
		//chsend, ok := n.(*ast.SendStmt)
		//if ok {
		//	fmt.Println(strconv.Itoa(fset.Position(chsend.Pos()).Line))
		//}
		//
		//// hladanie vsetkych prijati kanalom
		//chrecv, ok := n.(*ast.UnaryExpr)
		//if ok {
		//	if chrecv.Op.String() == "<-" {
		//		fmt.Println(strconv.Itoa(fset.Position(chrecv.Pos()).Line))
		//	}
		//}

		return true
	})

	//ff := ast.File{
	//	Name: ast.NewIdent("foo"),
	//	Decls: []ast.Decl{
	//		&ast.GenDecl{
	//			Tok: token.TYPE,
	//			Specs: []ast.Spec{
	//				&ast.TypeSpec{
	//					Name: ast.NewIdent("Bar"),
	//					Type: ast.NewIdent("uint32"),
	//				},
	//			},
	//		},
	//	},
	//}

	printTree(fset, node)

	//fmt.Println(repr(toString(fset, node)))
	//fmt.Println("All go funcs: ")
	//for _, goFunc := range goFuncs {
	//	fmt.Println(repr(goFunc))
	//}
}
