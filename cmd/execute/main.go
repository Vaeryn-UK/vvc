package main

import (
	"fmt"
	"github.com/vaeryn-uk/vvc/internal/core"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Must provide 1 input executable\n")
		os.Exit(1)
	}

	path := os.Args[1]

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		os.Stdout.WriteString("Could not read executable\n")
		os.Stdout.WriteString(fmt.Sprint(err))
		os.Exit(1)
	}

	if err := core.AutoConfigureMachine().Boot(core.BytesToWords(bytes)); err != nil {
		panic(err)
	}
}