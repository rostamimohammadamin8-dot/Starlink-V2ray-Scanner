package main

import (
	"fmt"
	"os"
)

func main() {
	mode := "collect"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	var err error
	switch mode {
	case "collect":
		err = collect()
	case "build":
		err = build()
	default:
		err = fmt.Errorf("unknown mode %q (use collect or build)", mode)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("Build Successful!")
}
