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

func DrawFunc(width uint, text string) (out [][]rune, height uint) {
	return box(width, text, rune('@'))
}

func DrawBox(width uint, text string) (out [][]rune, height uint) {
	return box(width, text, rune(':'))
}

func DrawFor(width uint, text string) (out [][]rune, height uint) {
	out, height = box(width, text, rune('*'))
	if 4 < width {
		out[0][0] = 'F'
		out[0][1] = 'O'
		out[0][2] = 'R'
		out[0][3] = ' '
	}
	return
}

func DrawIf(width uint, text string) (out [][]rune, height uint) {
	out, height = box(width, text, rune('#'))
	if 3 < width {
		out[0][0] = 'I'
		out[0][1] = 'F'
		out[0][2] = ' '
	}
	return
}

func ErrToOut(width uint, err error) (out string) {
	rs := []rune(fmt.Sprintf("%v", err))
	for iter := 0; iter < 1000; iter++ { // avoid infinite
		if int(width) < len(rs) {
			out += string(rs[:width]) + "\n"
			rs = rs[:width]
		} else {
			out += string(rs)
			break
		}
	}
	return
}

func Ascii(width uint, code string) (out string, err error) {
	if 1000 < width {
		width = 1000
	}
	// add package
	code = "package main\n" + code
	// gofmt code
	{
		var dat []byte
		var filename string
		var file *os.File
		if file, err = ioutil.TempFile("", "goast"); err != nil {
			out = ErrToOut(width, err)
			return
		}
		if _, err = file.WriteString(code); err != nil {
			out = ErrToOut(width, err)
			return
		}
		filename = file.Name()
		if err = file.Close(); err != nil {
			out = ErrToOut(width, err)
			return
		}
		if _, err = exec.Command("gofmt", "-w", filename).Output(); err != nil {
			out = ErrToOut(width, err)
			return
		}
		if dat, err = ioutil.ReadFile(filename); err != nil {
			out = ErrToOut(width, err)
			return
		}
		code = string(dat)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		out = ErrToOut(width, err)
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

func (v *Visitor) DrawNode(node ast.Node, dr func(width uint, text string) (out [][]rune, height uint)) {
	text := "undefined"
	if b, ok := node.(*ast.Ident); ok && b != nil {
		text = b.Name
	}
	if e, ok := node.(*ast.CallExpr); ok && e != nil {
		v.DrawNode(e.Fun, dr)
		return
	}
	if b, ok := node.(*ast.BasicLit); ok && b != nil {
		text = b.Value
	}
	out, _ := dr(v.width, text)
	view(&v.buf, out)
}

func lineLetter(buf io.Writer, width uint, letter rune) {
	rs := make([]rune, width)
	for i := range rs {
		rs[i] = ' '
	}
	index := width / 2
	if 1 < index {
		index--
	}
	if len(rs) == 0 {
		rs = []rune{' '}
		index = 0
	}
	rs[index] = letter
	fmt.Fprintf(buf, "%s\n", string(rs))
	fmt.Fprintf(buf, "%s\n", string(rs))
}

func line(buf io.Writer, width uint) {
	lineLetter(buf, width, '|')
}

func lineEmpty(buf io.Writer, width uint) {
	lineLetter(buf, width, ' ')
}

func block(width int, label string, list ast.Stmt) string {
	var r Visitor
	r.width = uint(width)
	astLabel := &ast.ExprStmt{X: &ast.BasicLit{Value: label}}
	if b, ok := list.(*ast.BlockStmt); ok && b != nil {
		b.List = append([]ast.Stmt{astLabel}, b.List...)
	}
	if list == nil {
		list = &ast.BlockStmt{
			List: []ast.Stmt{astLabel},
		}
	}
	r.Visit(list)
	return r.buf.String()
}

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {
	if f, ok := node.(*ast.File); ok && f != nil {
		for id, decl := range f.Decls {
			v.Visit(decl)
			if id != len(f.Decls)-1 {
				lineEmpty(&v.buf, v.width)
			}
			return
		}
	}
	if f, ok := node.(*ast.FuncDecl); ok && f != nil {
		out, _ := DrawFunc(v.width, f.Name.Name)
		view(&v.buf, out)
		line(&v.buf, v.width)
		for _, b := range f.Body.List {
			v.Visit(b)
			line(&v.buf, v.width)
		}
		out, _ = DrawFunc(v.width, fmt.Sprintf("End of %s", f.Name.Name))
		view(&v.buf, out)
	}
	if e, ok := node.(*ast.ExprStmt); ok && e != nil {
		v.Visit(e.X)
	}
	if i, ok := node.(*ast.ForStmt); ok && i != nil {
		v.DrawNode(i.Cond, DrawFor)
		// out, _ := DrawIf(v.width, "FOR")
		// view(&v.buf, out)
		// v.Visit(i.Body)
		// v.Visit(i.Cond)
		if v.width < 10 {
			return
		}
		leftWidth := 3
		left := " | "
		rightWidth := int(v.width) - leftWidth - 1
		right := block(rightWidth, "TRUE/ITERATE", i.Body)
		out := v.Merge(left, right)
		v.buf.WriteString(out)
		// end of if block
		v.DrawNode(&ast.BasicLit{Value: "End of for or iterate"}, DrawFor)
	}
	if block, ok := node.(*ast.BlockStmt); ok {
		line(&v.buf, v.width)
		for _, b := range block.List {
			v.Visit(b)
			line(&v.buf, v.width)
		}
		return
	}
	if i, ok := node.(*ast.IfStmt); ok && i != nil {
		v.DrawNode(i.Cond, DrawIf)
		// prepare blocks
		leftWidth := int(v.width)/2 - 1
		if i.Else == nil {
			leftWidth = int(v.width)*2/3 - 1
		}
		if len(i.Body.List) == 0 {
			leftWidth = int(v.width)*1/3 - 1
		}
		left := block(leftWidth, "TRUE", i.Body)
		rightWidth := int(v.width) - leftWidth - 1
		right := block(rightWidth, "FALSE", i.Else)
		out := v.Merge(left, right)
		v.buf.WriteString(out)
		// end of if block
		v.DrawNode(&ast.BasicLit{Value: "End of if"}, DrawIf)
	}
	if b, ok := node.(*ast.BasicLit); ok && b != nil {
		out, _ := DrawBox(v.width, b.Value)
		view(&v.buf, out)
	}
	if b, ok := node.(*ast.Ident); ok && b != nil {
		out, _ := DrawBox(v.width, b.Name)
		view(&v.buf, out)
	}
	if e, ok := node.(*ast.CallExpr); ok && e != nil {
		v.Visit(e.Fun)
	}
	return
}

func (v *Visitor) Merge(left, right string) string {
	ll := strings.Split(left, "\n")
	if ll[len(ll)-1] == "" {
		ll = ll[:len(ll)-1]
	}
	wl := len(ll[0])
	last := ll[len(ll)-1]
	lr := strings.Split(right, "\n")
	if lr[len(lr)-1] == "" {
		lr = lr[:len(lr)-1]
	}
	for i := 0; i < len(ll) && i < len(lr); i++ {
		ll[i] = ll[i] + " " + lr[i]
	}
	if len(ll) < len(lr) {
		rs := make([]rune, wl+1)
		for i := range rs {
			rs[i] = ' '
		}
		for i := len(ll); i < len(lr); i++ {
			ll = append(ll, last+" "+lr[i])
		}
	} else if len(lr) < len(ll) && 0 < len(lr) {
		for i := len(lr); i < len(ll); i++ {
			ll[i] = ll[i] + " " + lr[len(lr)-1]
		}
	}
	return strings.Join(ll, "\n") + "\n"
}

func view(buf io.Writer, out [][]rune) {
	for row := range out {
		for col := range out[row] {
			fmt.Fprintf(buf, "%s", string(out[row][col]))
		}
		fmt.Fprintf(buf, "\n")
	}
}
