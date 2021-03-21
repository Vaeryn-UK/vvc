package main

import (
	"bufio"
	"fmt"
	"github.com/vaeryn-uk/vvc/internal/compiler"
	"github.com/vaeryn-uk/vvc/internal/core"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Must provide 1 input bytecode file\n")
		os.Exit(1)
	}

	path := os.Args[1]

	os.Stderr.WriteString(fmt.Sprintf("Compiling file %s\n", path))

	file, err := os.Open(path)

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Could not open file %s (%s)\n", path, err))
		os.Exit(1)
	}

	c := compiler.NewCompiler()
	err, data := c.Compile(bufio.NewReader(file))

	if err != nil {
		os.Stderr.WriteString("Failed to compile\n")
		os.Stderr.WriteString(fmt.Sprintf("Reason: %s\n", err))
		os.Exit(1)
	}

	for _, word := range data {
		os.Stdout.Write(core.WordToBytes(word))
	}

	os.Stderr.WriteString("Done\n")
}