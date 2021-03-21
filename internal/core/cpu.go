package core

import (
	"errors"
	"fmt"
	"log"
)

type cpu struct {
	pc Address
	memory *memory
	flags CpuFlags
	requestStop bool
	haveSetPc bool
	registers []Word
}

type CpuFlags map[string]uint

func NewCpu(memory *memory, flags CpuFlags) *cpu {
	c := new(cpu)
	c.pc = 0
	c.memory = memory
	c.flags = flags

	c.registers = make([]Word, 8)

	return c
}

func (cpu *cpu) Wake() {
	for {
		cpu.cycle()

		if cpu.requestStop {
			return
		}
	}
}

func (cpu *cpu) cycle() {
	data := cpu.readNext()

	err, instruction := GetInstruction(Opcode(data))

	if err != nil {
		handleError(err)
	}

	args := make([]Word, 0)

	for i := uint8(0); i < instruction.ArgCount(); i++ {
		// Consume args.
		args = append(args, cpu.readNext())
	}

	cpu.execute(instruction, args)
}

func (cpu *cpu) execute(instruction Instruction, args []Word) {
	var argA, argB Word

	if len(args) > 0 {
		argA = args[0]
	}

	if len(args) > 1 {
		argB = args[1]
	}

	cpu.debug("Executing instruction %s with args (0x%x, 0x%x)", instruction.ref, argA, argB)

	cpu.haveSetPc = false

	switch instruction {
	case NOOP:
		// Nothing to do.
	case ADD:
		cpu.storeResult(cpu.getRegister(argA) + cpu.getRegister(argB), 0)
	case DEBUG:
		cpu.log("Register %d: %s", argA, cpu.getRegister(argA).ToString())
	case EXIT:
		cpu.requestStop = true
	case PUT:
		cpu.storeResult(argA, argB)
	case COPY:
		cpu.storeResult(cpu.getRegister(argA), argB)
	case COMPARE:
		a := cpu.getRegister(argA)
		b := cpu.getRegister(argB)

		var result Word

		if a < b {
			result = 1
		} else if a > b {
			result = 2
		} else {
			result = 0
		}

		cpu.storeResult(result, 0)
	case JEQ:
		if cpu.getRegister(0) == cpu.getRegister(argB) {
			cpu.jump(argA)
		}
	default:
		handleError(errors.New(fmt.Sprintf("Unhandled instruction %+v", instruction)))
	}
}

func (cpu *cpu) jump(pc Word) {
	cpu.debug("Jumping to PC %d", pc)
	cpu.haveSetPc = true
	cpu.pc = Address(pc)
}

func (cpu *cpu) storeResult(w Word, register Word) {
	cpu.assertRegister(register)

	cpu.registers[register] = w
}

func (cpu *cpu) assertRegister(number Word) {
	if uint(number) >= uint(len(cpu.registers)) {
		handleError(errors.New(fmt.Sprintf("Invalid register %d", number)))
	}
}

func (cpu *cpu) getRegister(number Word) Word {
	cpu.assertRegister(number)

	return cpu.registers[number]
}

func (cpu *cpu) debug(msg string, args ...interface{}) {
	debug, ok := cpu.flags["debug"]

	if ok && debug != 0 {
		cpu.log(msg, args...)
	}
}

func (cpu *cpu) log(msg string, args ...interface{}) {
	log.Printf("CPU: " + msg, args...)
}

func (cpu *cpu) readNext() Word {
	err, data := cpu.memory.Read(cpu.pc)
	cpu.pc++

	if err != nil {
		handleError(err)
	}

	return data
}

func handleError(err error) {
	panic(err)
}