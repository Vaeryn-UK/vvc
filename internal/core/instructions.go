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

// Noop does nothing.
var NOOP = Instruction{0, ArgTypeNone, ArgTypeNone, "NOP"}
// Adds number in register A to register B and stores the result in register 0.
var ADD = Instruction{1, ArgTypeRegister, ArgTypeRegister, "ADD"}
// Emits the current content of register A.
var DEBUG = Instruction{2, ArgTypeRegister, ArgTypeNone, "DBG"}
// Instructs the CPU to terminate.
var EXIT = Instruction{3, ArgTypeNone, ArgTypeNone, "EXT"}
// Loads data from addr A in to register B.
var LOAD = Instruction{4, ArgTypeData, ArgTypeRegister, "LOD"}
// Puts the literal A in to register B.
var PUT = Instruction{5, ArgTypeData, ArgTypeRegister, "PUT"}
// Shifts PC to memory A if register 0 is equal to register B.
var JEQ = Instruction{6, ArgTypeData, ArgTypeRegister, "JEQ"}
// Numeric comparison of register A and register B. Stores the result in register 0.
// rA = rB -> 0, rA < rB -> 1, rA > rB -> 2.
var COMPARE = Instruction{7, ArgTypeRegister, ArgTypeRegister, "CMP"}
// Copies the value in register A to register B.
var COPY = Instruction{8, ArgTypeRegister, ArgTypeRegister, "CPY"}

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