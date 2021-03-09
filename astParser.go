package KlaudiasGoTrace

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var goFuncs []string
var allFuncs []string
var fset *token.FileSet

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

//func newCallExpr(exp string) *ast.CallExpr {
//
//	funex, _ := parser.ParseExpr(exp)
//	expr := ast.CallExpr{
//		Fun:      funex,
//		Lparen:   0,
//		Args:     nil,
//		Ellipsis: 0,
//		Rparen:   0,
//	}
//	return &expr
//}

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

// vlozi hodnotu na dany index
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

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func indexOf(element interface{}, data []ast.Stmt) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func getChanNameAndVariable(line string) (string, string) {
	splitted := strings.Split(line, "<-")
	lastIndex := len(splitted) - 1
	chanName := strings.TrimSpace(splitted[lastIndex])
	var varName string
	if strings.Contains(line, ":=") {
		splitted2 := strings.Split(splitted[0], ":=")
		varName = strings.TrimSpace(splitted2[0])
	} else if strings.Contains(line, "=") {
		splitted2 := strings.Split(splitted[0], "=")
		varName = strings.TrimSpace(splitted2[0])
	}

	//fmt.Println("premenna:", varName, ", kanal:", chanName)
	return varName, chanName
}

// hladanie vsetkych posielani kanalom v deklaracii funkcie
func addSendToFuncDecl(funDecl *ast.FuncDecl) {
	ast.Inspect(funDecl, func(n ast.Node) bool {
		block, ok := n.(*ast.BlockStmt)
		if ok {
			var sends []int
			for i, stmt := range block.List {
				typeName := strings.TrimSpace(fmt.Sprintln(reflect.TypeOf(stmt)))
				if typeName == "*ast.SendStmt" {
					sends = append(sends, i)
				}
			}
			reverse(sends)

			for _, sendIndex := range sends {
				ast.Inspect(block.List[sendIndex], func(n ast.Node) bool {
					chanSend, ok := n.(*ast.SendStmt)
					if ok {
						//	fmt.Println(strconv.Itoa(fset.Position(chanSend.Pos()).Line), ", meno kanala:", chanName)
						chanName := fmt.Sprint(chanSend.Chan)
						value := fmt.Sprint(chanSend.Value)
						exprStr := fmt.Sprintf("KlaudiasGoTrace.SendToChannel(%s, %s)", value, chanName)
						expr3, _ := parser.ParseExpr(exprStr)
						stmt3 := ast.ExprStmt{X: expr3}
						block.List = insert(block.List, sendIndex, &stmt3)
					}
					return true
				})
			}
			return true
		}
		return true
	})
}

func addReceiveToFuncDecl(funDecl *ast.FuncDecl) {
	// hladanie vsetkych prijati kanalom
	//chrecv, ok := n.(*ast.UnaryExpr)
	//if ok {
	//	if chrecv.Op.String() == "<-" {
	//		fmt.Println(strconv.Itoa(fset.Position(chrecv.Pos()).Line))
	//	}
	//}
	astutil.Apply(funDecl, func(cursor *astutil.Cursor) bool {
		block, ok := cursor.Node().(*ast.BlockStmt)
		if ok {
			var recv [][]string
			astutil.Apply(block, func(cursor2 *astutil.Cursor) bool {
				chrecv, ok := cursor2.Node().(*ast.UnaryExpr)
				if ok && chrecv.Op.String() == "<-" {
					index := indexOf(cursor2.Parent(), block.List)
					if index >= 0 {
						values := []string{strconv.Itoa(index), toString(fset, cursor2.Parent())}
						recv = append(recv, values)
					}
				}
				return true
			}, nil)

			if len(recv) > 0 {
				reverse(recv)
				for i := range recv {
					recvIndex, _ := strconv.Atoi(recv[i][0])
					varName, chanName := getChanNameAndVariable(recv[i][1])
					withoutVar := false
					if varName == "" {
						withoutVar = true
						varName = "<-" + chanName
					}
					exprStr := fmt.Sprintf("KlaudiasGoTrace.ReceiveFromChannel(%s, %s)", varName, chanName)
					expr, _ := parser.ParseExpr(exprStr)
					stmt := ast.ExprStmt{X: expr}
					if withoutVar == true {
						block.List[recvIndex] = &stmt
					} else {
						block.List = insert(block.List, recvIndex+1, &stmt)
					}
				}
			}
		}
		return true
	}, nil)
}

func addExprToFuncDecl(funDecl *ast.FuncDecl, strExpr string, toStart bool) {
	expr, _ := parser.ParseExpr(strExpr)
	stmt := ast.ExprStmt{X: expr}
	if toStart == true {
		funDecl.Body.List = prepend(funDecl.Body.List, &stmt)
	} else {
		funDecl.Body.List = append(funDecl.Body.List, &stmt)
	}
}

func addParamToFuncDecl(funDecl *ast.FuncDecl, name string, typ string) {
	funcName := funDecl.Name.Name
	if funcName != "main" {
		field := newField(name, typ)
		funDecl.Type.Params.List = append(funDecl.Type.Params.List, field)
	}
}

func addExprToCall(callEx *ast.CallExpr) {
	funName := fmt.Sprint(callEx.Fun)
	if contains(goFuncs, funName) {
		expr, _ := parser.ParseExpr("KlaudiasGoTrace.GetGID()")
		callEx.Args = append(callEx.Args, expr)
	} else if contains(allFuncs, funName) {
		expr, _ := parser.ParseExpr("parentId")
		callEx.Args = append(callEx.Args, expr)
	}
}

func fullFillFuncArrays(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		// hladanie vsetkych deklaraácii funkcii, ktore vystupuju ako gorutiny
		fungo, ok := n.(*ast.GoStmt)
		if ok {
			name := getFuncName(toString(fset, fungo))
			goFuncs = append(goFuncs, name)
		}

		funDecl, ok := n.(*ast.FuncDecl)
		if ok {
			name := funDecl.Name.Name
			allFuncs = append(allFuncs, name)
		}
		return true
	})
}

func Parse(filePath string) {
	fset = token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fullFillFuncArrays(node)

	//fmt.Println("vsetky funkcie: " , allFuncs)
	//fmt.Println("vsetky go funkcie: ", goFuncs)

	ast.Inspect(node, func(n ast.Node) bool {
		funDecl, ok := n.(*ast.FuncDecl)
		if ok {
			addParamToFuncDecl(funDecl, "parentId", "uint64")
			addSendToFuncDecl(funDecl)
			addReceiveToFuncDecl(funDecl)
		}
		if ok && contains(goFuncs, funDecl.Name.Name) {
			//fmt.Println(strconv.Itoa(fset.Position(funDecl.Pos()).Line), funDecl.Name.Name, len(funDecl.Type.Params.List))
			addExprToFuncDecl(funDecl, "KlaudiasGoTrace.StartGoroutine(parentId)", true)
			addExprToFuncDecl(funDecl, "KlaudiasGoTrace.StopGoroutine()", false)
		}

		//hladanie vsetkych volani
		callEx, ok := n.(*ast.CallExpr)
		if ok {
			//fmt.Println(strconv.Itoa(fset.Position(callEx.Pos()).Line), callEx.Fun)
			addExprToCall(callEx)
		}

		// hladanie vsetkych prijati kanalom
		//chrecv, ok := n.(*ast.UnaryExpr)
		//if ok {
		//	if chrecv.Op.String() == "<-" {
		//		fmt.Println(strconv.Itoa(fset.Position(chrecv.Pos()).Line))
		//	}
		//}
		return true
	})

	printTree(fset, node)
}

//func main() {
//
//}

// vypisanie vsetkych parametrov funkcie
//for i := 0; i < len(funDecl.Type.Params.List); i++ {
//	fmt.Println(funDecl.Type.Params.List[i].Names[0].Name)
//}
