package flow

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Konstantin8105/compare"
)

func Test(t *testing.T) {
	tcs := []struct {
		isGraph bool
		name    string
		code    string
	}{
		////////
		{
			name: "empty",
			code: `
func Empty() {
}
		`},
		////////
		{
			name: "simple",
			code: `
func Simple() {
	"1"
	"2"
	"4345"
	"sdfd fsdfsad sda fad fa"
	"dsf asdfa;oieroi t[oprig fg ds ddsf akjl;dfk a;lsdkfa lkfg jsdfg"
	"Hello\nWorld"
	` + "`" + `
	Step
	One
	Two` + "`" + `
	` + "`" + `Step
	One
	Two` + "`" + `
}
        `,
		},
		////////
		{
			name: "if1",
			code: `
func OnlyIf() {
	if Find() {
		"CASE 1"
		"dfs kfja;lsdkfja;slkdfj a;lskdfj a;lskdjf a;lskdj fa;lskdjf asd"
		"sd faskdf asdfa dfa"
	}
}
		`,
		},
		////////
		{
			name: "for1",
			code: `
func OnlyFor() {
	for Conflict() {
		"Do some action"
	}
}
		`,
		},
		////////
		{
			name: "integration",
			code: `
func Action() {
	"Step 1"
	"Step 2"
	for Ecology() {
		"Step 3"
	}
	"Step 4"
	if Ecology2() {
		"Step 5"
	}
}
       `,
		},
		////////
		{
			name: "if2",
			code: `
func If2() {
	if Find() {
		"CASE 1"
		"dfs kfja;lsdkfja;slkdfj a;lskdfj a;lskdjf a;lskdj fa;lskdjf asd"
	} else {
		"CASE 2"
		"sdk fjasldkf jsdf lkasjdf ;laskdjf a"
		"d ldksjf a;ds fds dfsdfldkfj a;sldkfj asdfalskdfj"
	}
}
		`,
		},
		////////
		{
			name: "if3",
			code: `
func If3() {
	if Find() {
		"CASE 1"
		"dfs kfja;lsdkfja;slkdfj a;lskdfj a;lskdjf a;lskdj fa;lskdjf asd"
	} else {
		"CASE 2"
		"d ldksjf a;ds fds dfsdfldkfj a;sldkfj asdfalskdfj"
	}
}
		`,
		},
		////////
		{
			name: "if4",
			code: `
func If4() {
	if Find() {
	} else {
		"CASE 2"
		"d ldksjf a;ds fds dfsdfldkfj a;sldkfj asdfalskdfj"
	}
}
		`,
		},
		////////
		{
			name: "if5",
			code: `
func If5() {
	if "sd fasdfa kdjfal ksdjfla;ksdjf;laskdj f;alsdkj fasd fasdf asdf" {
	} else {
		"CASE 2"
		"d ldksjf a;ds fds dfsdfldkfj a;sldkfj asdfalskdfj"
	}
}
		`,
		},
		////////
		{
			name: "FuncDoc",
			code: `
func f0() {
}

// one line
func f1() {
}

// one line
// two line
func f2() {
}

/* one line
two line */
func f3() {
}
		`,
		},
		////////
		{
			name: "switch01",
			code: `
func switch01() {
	switch "General note" {
	case "State 1":
		"SQ1"
	case "State 2":
		"SW1"
		"SW2"
	case "State 3":
		"SR1"
	case "State 4":
		"sdfsd fkjsdlk fjasdlk fjdslkf jasd;lkfj as;dlkf jasdl;fk ja;sdlkfj asdf asdf asdfasd"
	case "State 5":
		"sdfsd fkjsdlk fs"
	case "df seroirpwoer iwqepori qw[eori q[pweoirq[eoir q[wepor iqwe":
		"SD"
	default:
		"DF1"
	}
}
		`,
		},
		////////
		{
			name:    "graph1",
			isGraph: true,
			code: `
			"   1" > "  2  " > ` + "`" + `3 "d"f` + "`" + `   < "4" < "5" < "6

			" ` + " " + `
			`,
		},
		////////
		{
			name:    "graph2",
			isGraph: true,
			code:    ` "1" > `+"`" + `2
			sd fsd f 
			` + "`" + ` > "3"`,
		},
		////////
		{
			name:    "funcs",
			isGraph: false,
			code:    `
func c1() {
	c2()
	"a1"
	c3()
	"a2"
	c4()
	c5()
	c6("dsfssdfsdfs")
}`,
		},
		////////
	}
	for _, tc := range tcs {
		for _, width := range []uint{5, 10, 15, 20, 31, 40} {
			filename := fmt.Sprintf("testdata/%s_W%03d", tc.name, width)
			t.Run(filename, func(t *testing.T) {
				debug = testing.Verbose()
				var out string
				var err error
				if tc.isGraph {
					out, err = Graph(width, tc.code)
				} else {
					out, err = Ascii(width, tc.code)
				}
				if err != nil {
					t.Fatalf("Error:\n%s\n", err)
					return
				}
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "%s\n%s", tc.code, out)
				compare.Test(t, filename, buf.Bytes())
				debug = false
			})
		}
	}

}

func ExampleDrawText() {
	var width uint = 10
	out, height := DrawText(width, "Long lorem porem text")
	err := draw(width, height, out)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		return
	}
	// Output:
	// |LONG LOREM|
	// | POREM TEX|
	// |T         |
}

func ExampleDrawBox() {
	var width uint = 10
	out, height := DrawBox(width, "Long lorem porem text")
	err := draw(width, height, out)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		return
	}
	// Output:
	// |::::::::::|
	// |: LONG L :|
	// |: OREM P :|
	// |: OREM T :|
	// |: EXT    :|
	// |::::::::::|
}

func ExampleDrawIf() {
	var width uint = 10
	out, height := DrawIf(width, "Long lorem porem text")
	err := draw(width, height, out)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		return
	}
	// Output:
	// |IF #######|
	// |# LONG L #|
	// |# OREM P #|
	// |# OREM T #|
	// |# EXT    #|
	// |##########|
}

func draw(width, height uint, out [][]rune) error {
	if len(out) != int(height) {
		return fmt.Errorf(
			"height is not same: %d != %d",
			len(out),
			height,
		)
	}
	for row := range out {
		if len(out[row]) != int(width) {
			return fmt.Errorf(
				"width is not same: %d != %d",
				len(out[row]),
				width,
			)
		}
	}
	for row := range out {
		fmt.Fprintf(os.Stdout, "|")
		for col := range out[row] {
			fmt.Fprintf(os.Stdout, "%s", string(out[row][col]))
		}
		fmt.Fprintf(os.Stdout, "|\n")
	}
	return nil
}
