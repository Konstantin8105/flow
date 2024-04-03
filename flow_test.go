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
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "%s", out)
	// Output:
}
