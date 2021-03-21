package core

import (
	"encoding/binary"
	"fmt"
)

type Word uint32

const bytesInWord = 4

func WordToBytes(w Word) []byte {
	bytes := make([]byte, bytesInWord)
	binary.LittleEndian.PutUint32(bytes, uint32(w))
	return bytes
}

func BytesToWords(bytes []byte) Words {
	words := make(Words, 0)

	for i := 0; i < len(bytes); i += bytesInWord {
		words = append(words, Word(binary.LittleEndian.Uint32(bytes[i:i+bytesInWord])))
	}

	return words
}

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