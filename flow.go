package flow

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Konstantin8105/tf"
)

var debug bool // debug info

func box(width uint, text string, border rune) (out [][]rune, height uint) {
	{ // cleaning text
		text = strings.ReplaceAll(text, "\t", "")
		text = strings.TrimPrefix(text, "\"")
		text = strings.TrimSuffix(text, "\"")
		text = strings.TrimPrefix(text, "`")
		text = strings.TrimSuffix(text, "`")
		lines := strings.Split(text, "\n")
		for i := range lines {
			lines[i] = strings.TrimSpace(lines[i])
		}
	again:
		for i := range lines {
			if lines[i] != "" {
				continue
			}
			lines = append(lines[:i], lines[i+1:]...)
			goto again
		}
		text = strings.Join(lines, "\n")
	}
	if width < 4 {
		line := make([]rune, width)
		for col := range line {
			line[col] = border
		}
		return [][]rune{line}, 1
	}
	out, height = DrawText(width-4, text)
	for row := range out {
		out[row] = append([]rune{border, ' '}, append(out[row], ' ', border)...)
	}
	out = append(make([][]rune, 1), append(out, make([][]rune, 1)...)...)
	out[0] = make([]rune, width)
	out[len(out)-1] = make([]rune, width)
	for col := range out[0] {
		out[0][col] = border
		out[len(out)-1][col] = border
	}
	height += 2
	return
}

func DrawText(width uint, text string) (out [][]rune, height uint) {
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

func DrawBox(width uint, text string) (out [][]rune, height uint) {
	return box(width, text, rune('*'))
}

func DrawIf(width uint, text string) (out [][]rune, height uint) {
	out, height = box(width, text, rune('#'))
	// 	out[0][0] = 'I'
	// 	out[0][1] = 'F'
	// 	out[0][2] = ' '
	return
}

func Ascii(width uint, code string) (out string, err error) {
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
	v := Visitor{width: width}
	ast.Walk(&v, f)
	out = v.buf.String()
	if debug {
		var buf bytes.Buffer
		ast.Fprint(&buf, fset, f, ast.NotNilFilter)
		result := buf.String()
		fmt.Println(result)
	}
	return
}

type Visitor struct {
	width uint
	buf   bytes.Buffer
}

var tab uint

func line(buf io.Writer, width uint) {
	rs := make([]rune, width)
	for i := range rs {
		rs[i] = ' '
	}
	index := width / 2
	if 1 < index {
		index--
	}
	rs[index] = '|'
	fmt.Fprintf(buf, "%s\n", string(rs))
	fmt.Fprintf(buf, "%s\n", string(rs))
}

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {
	var width uint = v.width - tab*3
	if f, ok := node.(*ast.File); ok && f != nil {
		for _, decl := range f.Decls {
			fmt.Fprintln(&v.buf, ">")
			v.Visit(decl)
			fmt.Fprintln(&v.buf, "<")
			return
		}
	}
	if f, ok := node.(*ast.FuncDecl); ok && f != nil {
		for _, b := range f.Body.List {
			v.Visit(b)
			line(&v.buf, width)
		}
	}
	if e, ok := node.(*ast.ExprStmt); ok && e != nil {
		v.Visit(e.X)
	}
	if i, ok := node.(*ast.ForStmt); ok && i != nil {
		out, _ := DrawIf(width, "FOR")
		view(&v.buf, out)
		line(&v.buf, width)
		tab++
		for _, b := range i.Body.List {
			v.Visit(b)
			line(&v.buf, width)
		}
		tab--
		v.Visit(i.Cond)
	}
	if i, ok := node.(*ast.IfStmt); ok && i != nil {
		v.Visit(i.Cond)
		line(&v.buf, width)
		tab++
		for _, b := range i.Body.List {
			v.Visit(b)
			line(&v.buf, width)
		}
		tab--
		out, _ := DrawIf(width, "Enf of if")
		view(&v.buf, out)
	}
	if b, ok := node.(*ast.BasicLit); ok && b != nil {
		out, _ := DrawBox(width, b.Value)
		view(&v.buf, out)
	}
	if b, ok := node.(*ast.Ident); ok && b != nil {
		out, _ := DrawIf(width, b.Name)
		view(&v.buf, out)
	}
	if e, ok := node.(*ast.CallExpr); ok && e != nil {
		v.Visit(e.Fun)
	}
	return
}

func view(buf io.Writer, out [][]rune) {
	for row := range out {
		for col := range out[row] {
			fmt.Fprintf(buf, "%s", string(out[row][col]))
		}
		fmt.Fprintf(buf, "\n")
	}
}
