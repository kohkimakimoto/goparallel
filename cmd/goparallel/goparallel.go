package main

import (
	"fmt"
	"github.com/kohkimakimoto/goparallel/goparallel"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "goparallel error: %v\n", err)
			os.Exit(1)
		}
	}()

	if err := goparallel.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "goparallel error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
