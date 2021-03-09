package KlaudiasGoTrace

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var goFuncs []string
var allFuncs []string
var fset *token.FileSet
var root *ast.File
var fileName string

func toString(fset *token.FileSet, file interface{}) string {
	var buffer bytes.Buffer
	printer.Fprint(&buffer, fset, file)
	return prettify(buffer.String())
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
	return varName, chanName
}

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

func addAssignStmt(funcDecl *ast.FuncDecl, left string, typ token.Token, right string) {
	leftExpr, _ := parser.ParseExpr(left)
	rightExpr, _ := parser.ParseExpr(right)
	assign := ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{leftExpr},
		Rhs: []ast.Expr{rightExpr},
	}
	funcDecl.Body.List = prepend(funcDecl.Body.List, &assign)
}

func addStartStopToMain(funDecl *ast.FuncDecl) {
	addExprToFuncDecl(funDecl, "KlaudiasGoTrace.StartTrace()", true)
	addExprToFuncDecl(funDecl, "KlaudiasGoTrace.EndTrace()", false)
}

func addImport(importString string) {
	ast.Inspect(root, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if ok {
			if genDecl.Tok == token.IMPORT {
				iSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote(importString)}}
				genDecl.Specs = append(genDecl.Specs, iSpec)
			}
		}
		return true
	})
}

func prettify(uglyStr string) string {
	prettyStr := strings.ReplaceAll(uglyStr, "KlaudiasGoTrace.\n\t\t", "KlaudiasGoTrace.")
	prettyStr = strings.ReplaceAll(prettyStr, "KlaudiasGoTrace.\t", "KlaudiasGoTrace.")
	prettyStr = strings.ReplaceAll(prettyStr, ",\n\t\t\t)", ",)")
	prettyStr = strings.ReplaceAll(prettyStr, ",\n\t\t)", ",)")
	prettyStr = strings.ReplaceAll(prettyStr, ",\n\n", ",")
	prettyStr = strings.ReplaceAll(prettyStr, ",\n\n\t\t\t\t", ", ")
	prettyStr = strings.ReplaceAll(prettyStr, "KlaudiasGoTrace.StartGoroutine(\n\t\t", "KlaudiasGoTrace.StartGoroutine(")
	prettyStr = strings.ReplaceAll(prettyStr, ",\t\t\t\t", ", ")
	prettyStr = strings.ReplaceAll(prettyStr, "KlaudiasGoTrace.GetGID(),\n", "KlaudiasGoTrace.GetGID(),")
	return prettyStr
}

func writeToFile(filename string, data string) bool {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func createFileFromAST(filename string, data string) string {
	fileVersion := 0
	formatPostfix := "Parsed_.go"
	postfix := "Parsed.go"
	fileName := strings.ReplaceAll(filename, ".go", postfix)

	for !writeToFile(fileName, data) {
		fileVersion += 1
		postfix = strings.ReplaceAll(formatPostfix, "_", strconv.Itoa(fileVersion))
		fileName = strings.ReplaceAll(filename, ".go", postfix)
	}

	return fileName
}

func Parse(filePath string) string {
	fset = token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var results []string

	results = strings.Split(filePath, "\\")
	fileName = results[len(results)-1]

	root = node

	fullFillFuncArrays(node)
	addImport("KlaudiasGoTrace/KlaudiasGoTrace")

	ast.Inspect(node, func(n ast.Node) bool {
		funDecl, ok := n.(*ast.FuncDecl)
		if ok {
			funcName := funDecl.Name.Name
			addParamToFuncDecl(funDecl, "parentId", "uint64")
			addSendToFuncDecl(funDecl)
			addReceiveToFuncDecl(funDecl)

			if contains(goFuncs, funcName) {
				addExprToFuncDecl(funDecl, "KlaudiasGoTrace.StartGoroutine(parentId)", true)
				addExprToFuncDecl(funDecl, "KlaudiasGoTrace.StopGoroutine()", false)
			}

			if funcName == "main" {
				addAssignStmt(funDecl, "parentId", token.INT, "uint64(0)")
				addStartStopToMain(funDecl)
			}
		}

		callEx, ok := n.(*ast.CallExpr)
		if ok {
			addExprToCall(callEx)
		}

		return true
	})
	//printTree(fset, node)

	return createFileFromAST(fileName, toString(fset, node))
}
