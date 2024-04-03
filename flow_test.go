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

func ExampleConvert() {
	var width uint = 10
	out, height := convert(10, "Long lorem porem text")
	if len(out) != int(height) {
		fmt.Fprintf(os.Stdout,
			"height is not same: %d != %d",
			len(out),
			height,
		)
		return
	}
	for row := range out {
		if len(out[row]) != int(width) {
			fmt.Fprintf(os.Stdout,
				"width is not same: %d != %d",
				len(out[row]),
				width,
			)
			return
		}
		fmt.Fprintf(os.Stdout, "|")
		for col := range out[row] {
			fmt.Fprintf(os.Stdout, "%s", string(out[row][col]))
		}
		fmt.Fprintf(os.Stdout, "|\n")
	}
	// Output:
	// |Long lorem|
	// | porem tex|
	// |t         |
}
