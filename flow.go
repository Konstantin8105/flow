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

var (
	RuneFunc     = '@'
	RuneBox      = ':'
	RuneFor      = '*'
	RuneIf       = '#'
	RuneSwitch   = '$'
	RuneDown     = 'V'
	RuneUp       = '^'
	RuneVertical = '|'
)

func box(width uint, text string, border rune) (out [][]rune, height uint) {
	text = strings.ToUpper(text)
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
	text = strings.ToUpper(text)
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
	return box(width, text, RuneFunc)
}

func DrawBox(width uint, text string) (out [][]rune, height uint) {
	return box(width, text, RuneBox)
}

func DrawFor(width uint, text string) (out [][]rune, height uint) {
	out, height = box(width, text, RuneFor)
	if 4 < width {
		out[0][0] = 'F'
		out[0][1] = 'O'
		out[0][2] = 'R'
		out[0][3] = ' '
	}
	return
}

func DrawSwitch(width uint, text string) (out [][]rune, height uint) {
	out, height = box(width, text, RuneSwitch)
	if 7 < width {
		out[0][0] = 'S'
		out[0][1] = 'W'
		out[0][2] = 'I'
		out[0][3] = 'T'
		out[0][4] = 'C'
		out[0][5] = 'H'
		out[0][6] = ' '
	}
	return
}

func DrawIf(width uint, text string) (out [][]rune, height uint) {
	out, height = box(width, text, RuneIf)
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

func Graph(width uint, code string) (out string, err error) {
	rs := []rune(code)
	var names []string
	var inC bool
	var inT bool
	var links []bool
	for i := range rs {
		if !inT && rs[i] == '"' {
			if !inC {
				names = append(names, "")
			}
			inC = !inC
			if inC {
				continue
			}
		}
		if !inC && rs[i] == '`' {
			if !inT {
				names = append(names, "")
			}
			inT = !inT
			if inT {
				continue
			}
		}
		if inC || inT {
			names[len(names)-1] += string(rs[i])
		} else if rs[i] != ' ' && rs[i] != '"' && rs[i] != '`' && !inC && !inT {
			if rs[i] == '>' {
				links = append(links, true)
				continue
			} else if rs[i] == '<' {
				links = append(links, false)
				continue
			}
			if debug {
				fmt.Fprintf(os.Stdout, "%d %s\n", i, string(rs[i]))
			}
		}
	}
	var buf bytes.Buffer
	for i := range names {
		out, _ := DrawBox(width, names[i])
		view(&buf, out)
		if i != len(names)-1 {
			if i == len(links) {
				err = fmt.Errorf("link error")
				continue
			}
			if links[i] {
				line(&buf, width)
				lineLetter(&buf, width, RuneDown)
			} else {
				lineLetter(&buf, width, RuneUp)
				line(&buf, width)
			}
		}
	}
	out = buf.String()
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
			err = fmt.Errorf("TempFile: %v", err)
			out = ErrToOut(width, err)
			return
		}
		if _, err = file.WriteString(code); err != nil {
			err = fmt.Errorf("WriteString: %v", err)
			out = ErrToOut(width, err)
			return
		}
		filename = file.Name()
		if err = file.Close(); err != nil {
			err = fmt.Errorf("file Close: %v", err)
			out = ErrToOut(width, err)
			return
		}
		if _, err = exec.Command("gofmt", "-w", filename).Output(); err != nil {
			err = fmt.Errorf("gofmt: %v", err)
			out = ErrToOut(width, err)
			return
		}
		if dat, err = ioutil.ReadFile(filename); err != nil {
			err = fmt.Errorf("read file: %v", err)
			out = ErrToOut(width, err)
			return
		}
		code = string(dat)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		err = fmt.Errorf("parse file: %v", err)
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
	text := fmt.Sprintf("undefined in DrawNode: %T", node)
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
	if es, ok := node.(*ast.ExprStmt); ok && es != nil {
		v.DrawNode(es.X, dr)
		return
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
	// fmt.Fprintf(buf, "%s\n", string(rs))
}

func line(buf io.Writer, width uint) {
	lineLetter(buf, width, RuneVertical)
}

func lineEmpty(buf io.Writer, width uint) {
	lineLetter(buf, width, ' ')
}

func block(width int, label string, list ast.Node) string {
	var r Visitor
	r.width = uint(width)
	astLabel := &ast.ExprStmt{X: &ast.BasicLit{Value: label}}
	if b, ok := list.(*ast.BlockStmt); ok && b != nil && label != "" {
		b.List = append([]ast.Stmt{astLabel}, b.List...)
	}
	if list == nil {
		list = &ast.BlockStmt{
			List: []ast.Stmt{astLabel},
		}
	}
	// if begin {
	line(&r.buf, r.width)
	// }
	r.Visit(list)
	// line(&r.buf, r.width)
	return r.buf.String()
}

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.File:
		for id, decl := range n.Decls {
			v.Visit(decl)
			if id != len(n.Decls)-1 {
				lineEmpty(&v.buf, v.width)
			}
		}
	case *ast.FuncDecl:
		docs := getDocs(n.Doc)
		out, _ := DrawFunc(v.width, n.Name.Name+"\n"+docs)
		view(&v.buf, out)
		line(&v.buf, v.width)
		v.Visit(n.Body)
		out, _ = DrawFunc(v.width, fmt.Sprintf("End of %s", n.Name.Name))
		view(&v.buf, out)
	case *ast.ExprStmt:
		v.Visit(n.X)
	case *ast.ForStmt:
		v.DrawNode(n.Cond, DrawFor)
		if v.width < 10 {
			return
		}
		leftWidth := 3
		left := " " + string(RuneVertical) + " "
		rightWidth := int(v.width) - leftWidth - 1
		right := block(rightWidth, "TRUE/ITERATE", n.Body)
		out := v.Merge(left, right)
		v.buf.WriteString(out)
		// end of if block
		v.DrawNode(&ast.BasicLit{Value: "End of for or iterate"}, DrawFor)
	case *ast.BlockStmt:
		for _, b := range n.List {
			v.Visit(b)
			line(&v.buf, v.width)
		}
	case *ast.SwitchStmt:
		v.DrawNode(n.Tag, DrawSwitch)
		line(&v.buf, v.width)
		for _, b := range n.Body.List {
			v.Visit(b)
		}
		v.DrawNode(&ast.BasicLit{Value: "End of switch"}, DrawSwitch)
	case *ast.CaseClause:
		if len(n.List) == 1 && len(n.Body) == 1 {
			leftWidth := uint(v.width) / 2
			rightWidth := uint(v.width) - leftWidth - 1
			left := Visitor{width: leftWidth}
			left.DrawNode(n.List[0], DrawIf)
			{
				rs := make([]rune, leftWidth+1)
				for i := range rs {
					rs[i] = ' '
				}
				rs[1] = RuneVertical
				rs[len(rs)-1] = '\n'
				fmt.Fprintf(&left.buf, string(rs))
			}
			right := Visitor{width: rightWidth}
			right.DrawNode(n.Body[0], DrawBox)
			{
				rs := make([]rune, leftWidth+1)
				for i := range rs {
					rs[i] = ' '
				}
				rs[len(rs)-1] = '\n'
				fmt.Fprintf(&right.buf, string(rs))
			}
			out := v.Merge(left.buf.String(), right.buf.String())
			v.buf.WriteString(out)
			break
		}
		for i := range n.List {
			v.DrawNode(n.List[i], DrawIf)
		}
		if len(n.List) == 0 {
			v.DrawNode(&ast.BasicLit{Value: "Default case of switch"}, DrawIf)
		}
		left := " " + string(RuneVertical) + " "
		rightWidth := int(v.width) - len([]rune(left)) - 1
		right := block(rightWidth, "", &ast.BlockStmt{List: n.Body})
		out := v.Merge(left, right)
		v.buf.WriteString(out)
	case *ast.IfStmt:
		v.DrawNode(n.Cond, DrawIf)
		// prepare blocks
		leftWidth := int(v.width)/2 - 1
		if n.Else == nil {
			rightWidth := min(int(v.width)*1/3, 10)
			if rightWidth < 3 {
				rightWidth = 3
			}
			leftWidth = int(v.width) - rightWidth
		}
		if len(n.Body.List) == 0 {
			leftWidth = min(int(v.width)*1/3, 8)
		}
		left := block(leftWidth, "TRUE", n.Body)
		rightWidth := int(v.width) - leftWidth - 1
		right := block(rightWidth, "FALSE", n.Else)
		out := v.Merge(left, right)
		v.buf.WriteString(out)
		// end of if block
		v.DrawNode(&ast.BasicLit{Value: "End of if"}, DrawIf)
	case *ast.BasicLit:
		out, _ := DrawBox(v.width, n.Value)
		view(&v.buf, out)
	case *ast.Ident:
		out, _ := DrawBox(v.width, n.Name)
		view(&v.buf, out)
	case *ast.CallExpr:
		v.Visit(n.Fun)
	default:
		out, _ := DrawBox(v.width, fmt.Sprintf("UNDEFINED: %T", n))
		view(&v.buf, out)
	}
	return
}

func (v *Visitor) Merge(left, right string) string {
	ll := strings.Split(left, "\n")
	if ll[len(ll)-1] == "" {
		ll = ll[:len(ll)-1]
	}
	wl := len(ll[0])
	lr := strings.Split(right, "\n")
	if lr[len(lr)-1] == "" {
		lr = lr[:len(lr)-1]
	}
	// last lines
	lastLeft := ll[len(ll)-1]
	if !strings.Contains(lastLeft, string(RuneVertical)) {
		rs := make([]rune, len([]rune(lastLeft)))
		for i := range rs {
			rs[i] = ' '
		}
		lastLeft = string(rs)
	}
	lastRight := lr[len(lr)-1]
	if !strings.Contains(lastRight, string(RuneVertical)) {
		rs := make([]rune, len([]rune(lastRight)))
		for i := range rs {
			rs[i] = ' '
		}
		lastRight = string(rs)
	}
	//
	for i := 0; i < len(ll) && i < len(lr); i++ {
		ll[i] = ll[i] + " " + lr[i]
	}
	if len(ll) < len(lr) {
		rs := make([]rune, wl+1)
		for i := range rs {
			rs[i] = ' '
		}
		for i := len(ll); i < len(lr); i++ {
			ll = append(ll, lastLeft+" "+lr[i])
		}
	} else if len(lr) < len(ll) && 0 < len(lr) {
		for i := len(lr); i < len(ll); i++ {
			ll[i] = ll[i] + " " + lastRight
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

func getDocs(d *ast.CommentGroup) (out string) {
	if d == nil {
		return
	}
	for i := range d.List {
		out += d.List[i].Text
		if i != len(d.List)-1 {
			out += "\n"
		}
	}
	return
}
