package core

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