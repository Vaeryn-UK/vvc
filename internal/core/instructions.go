package core

import (
	"errors"
	"fmt"
)

type Opcode uint

type Instruction struct {
	opcode Opcode
	args uint8
	ref string
}

var NOOP = Instruction{0, 0, "NOOP"}
var ADD = Instruction{1, 2, "ADD"}
var DEBUG = Instruction{2, 0, "DEBUG"}
var EXIT = Instruction{3, 0, "EXIT"}

var instructionsByRef = map[string]Instruction{
	NOOP.ref: NOOP,
	ADD.ref: ADD,
	DEBUG.ref: DEBUG,
	EXIT.ref: EXIT,
}

func LookupInstruction(ref string) (error, Instruction) {
	i, ok := instructionsByRef[ref]

	if !ok {
		return errors.New("Invalid instruction: " + ref), i
	}

	return nil, i
}

func GetInstruction(opcode Opcode) (error, Instruction) {
	var i Instruction

	// TODO: speed this up.
	for _, instruction := range instructionsByRef {
		if instruction.opcode == opcode {
			return nil, instruction
		}
	}

	return errors.New(fmt.Sprintf("Unknown instruction 0x%x", opcode)), i
}