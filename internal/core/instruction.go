package core

import (
	"errors"
	"fmt"
)

type Opcode uint

type ArgType uint

const (
	ArgTypeNone ArgType = iota
	ArgTypeRegister
	ArgTypeData
)

type Instruction struct {
	opcode Opcode
	A      ArgType
	B      ArgType
	ref    string
}

func (i *Instruction) Opcode() Opcode {
	return i.opcode
}

func (i *Instruction) ArgCount() uint8 {
	count := uint8(0)

	if i.A != ArgTypeNone {
		count++
	}

	if i.B != ArgTypeNone {
		count++
	}

	return count
}

var instructionsByRef = map[string]Instruction{
	NOOP.ref: NOOP,
	ADD.ref: ADD,
	DEBUG.ref: DEBUG,
	EXIT.ref: EXIT,
	LOAD.ref: LOAD,
	PUT.ref: PUT,
	JEQ.ref: JEQ,
	COMPARE.ref: COMPARE,
	COPY.ref: COPY,
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