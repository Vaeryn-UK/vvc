package core

import (
	"fmt"
)

type Word uint32

type Address Word

func IsValidWord(val interface{}) bool {
	return true
}

type Words []Word

func (words Words) ToString() string {
	result := fmt.Sprintf("Word[] size: %d.\n", len(words))

	for _, w := range words {
		result += fmt.Sprintf("0x%x", w) + "\n"
	}

	return result
}