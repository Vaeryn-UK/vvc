package main

import (
	"bufio"
	"fmt"
	"github.com/vaeryn-uk/vvc/internal/core"
	"log"
	"os"
)

func main() {
	path := "programs/add.vvb"

	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	err, program := core.NewCompiler().Compile(bufio.NewReader(file))

	if err != nil {
		log.Fatal(err)
	}

	cpuFlags := make(map[string]int)

	if os.Getenv("VVC_DEBUG") == "1" {
		cpuFlags["debug"] = 1
	}

	fmt.Printf("%s compiled.", path)

	machine := core.NewMachine(256, cpuFlags)

	if err := machine.Boot(program); err != nil {
		panic(err)
	}
}