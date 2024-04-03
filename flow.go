package flow

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Konstantin8105/tf"
)

func convert(width uint, text string) (out [][]rune, height uint) {
	width += 1
	var b tf.Buffer
	var t tf.TextField
	t.SetWidth(width)
	t.SetText([]rune(text))
	t.Render(b.Drawer, nil)
	for row := range b {
		if int(width) <= len(b[row]) {
			continue
		}
		size := int(width) - len(b[row])
		addition := make([]rune, size)
		for i := range addition {
			addition[i] = ' '
		}
		b[row] = append(b[row], addition...)
	}
	for row := range b {
		b[row] = b[row][:len(b[row])-1]
	}
	out = b
	height = uint(len(b))
	return
}

func Ascii(code string) (out string, err error) {
	// add package
	code = "package main\n" + code
	// gofmt code
	{
		var dat []byte
		var filename string
		var file *os.File
		if file, err = ioutil.TempFile("", "goast"); err != nil {
			return
		}
		if _, err = file.WriteString(code); err != nil {
			return
		}
		filename = file.Name()
		if err = file.Close(); err != nil {
			return
		}
		if _, err = exec.Command("gofmt", "-w", filename).Output(); err != nil {
			return
		}
		if dat, err = ioutil.ReadFile(filename); err != nil {
			return
		}
		code = string(dat)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	var v Visitor
	ast.Walk(&v, f)
	// 	{
	// 		var buf bytes.Buffer
	// 		ast.Fprint(&buf, fset, f, ast.NotNilFilter)
	// 		result := buf.String()
	// 		fmt.Println(result)
	// 	}
	return
}

type Visitor struct{}

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {
	if f, ok := node.(*ast.File); ok && f != nil {
		for _, decl := range f.Decls {
			fmt.Println(">")
			v.Visit(decl)
			fmt.Println("<")
			return
		}
	}
	if f, ok := node.(*ast.FuncDecl); ok && f != nil {
		for _, b := range f.Body.List {
			v.Visit(b)
		}
	}
	if e, ok := node.(*ast.ExprStmt); ok && e != nil {
		v.Visit(e.X)
	}
	if b, ok := node.(*ast.BasicLit); ok && b != nil {
		fmt.Println(b.Value)
	}
	return
}
