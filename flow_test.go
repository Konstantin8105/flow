package flow

import (
	"fmt"
	"os"
)

func Example() {
	code := `
func Action() {
	"Step 1"
	"Step 2"
	if Ecology() {
		"Step 3"
	}
	"Step 4"
}`
	out, err := Ascii(code)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error:\n%s\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "%s", out)
	// Output:
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
	// |IF #######|
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
		fmt.Fprintf(os.Stdout, "|")
		for col := range out[row] {
			fmt.Fprintf(os.Stdout, "%s", string(out[row][col]))
		}
		fmt.Fprintf(os.Stdout, "|\n")
	}
	return nil
}
