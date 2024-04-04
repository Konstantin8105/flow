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
		name string
		code string
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
func If3() {
	if Find() {
	} else {
		"CASE 2"
		"d ldksjf a;ds fds dfsdfldkfj a;sldkfj asdfalskdfj"
	}
}
		`,
		},
	}
	for _, tc := range tcs {
		for _, width := range []uint{5, 10, 15, 20, 40} {
			filename := fmt.Sprintf("testdata/%s_W%03d", tc.name, width)
			t.Run(filename, func(t *testing.T) {
				debug = testing.Verbose()
				out, err := Ascii(width, tc.code)
				if err != nil {
					t.Fatalf("Error:\n%s\n", err)
					return
				}
				var buf bytes.Buffer
				// for row := range out {
				// 	for col := range out[row] {
				// 		fmt.Fprintf(&buf, "%s", string(out[row][col]))
				// 	}
				// 	fmt.Fprintf(&buf, "\n")
				// }
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
	// |Long lorem|
	// | porem tex|
	// |t         |
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
	// |**********|
	// |* Long l *|
	// |* orem p *|
	// |* orem t *|
	// |* ext    *|
	// |**********|
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
	// |##########|
	// |# Long l #|
	// |# orem p #|
	// |# orem t #|
	// |# ext    #|
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
