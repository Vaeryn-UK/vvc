package core

import "log"

type cpu struct {
	pc Address
	memory *memory
	lastResult Word
	flags map[string]int
	requestStop bool
}

func NewCpu(memory *memory, flags map[string]int) *cpu {
	c := new(cpu)
	c.pc = 0
	c.memory = memory
	c.flags = flags
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

	for i := uint8(0); i < instruction.args; i++ {
		// Consume args.
		args = append(args, cpu.readNext())
	}

	cpu.execute(instruction, args)
}

func (cpu *cpu) execute(instruction Instruction, args []Word) {
	var result Word

	cpu.debug("Executing instruction %v", instruction)

	switch instruction {
	case NOOP:
		// Nothing to do.
	case ADD:
		result = args[0] + args[1]
	case DEBUG:
		cpu.log("Last result: %v", cpu.lastResult)
	case EXIT:
		cpu.requestStop = true
	}

	cpu.lastResult = result
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