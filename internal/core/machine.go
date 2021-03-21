package core

import "os"

type Machine struct {
	memory *memory
	cpu *cpu
}

func NewMachine(memSize uint, cpuFlags map[string]int) *Machine {
	m := NewMemory(memSize)

	return &Machine{memory: m, cpu: NewCpu(m, cpuFlags)}
}

func AutoConfigureMachine() *Machine {
	cpuFlags := make(map[string]int)

	if os.Getenv("VVC_DEBUG") == "1" {
		cpuFlags["debug"] = 1
	}

	memSize := uint(256)

	return NewMachine(memSize, cpuFlags)
}

func (m *Machine) Boot(program Words) error {
	if err := m.memory.Load(0, program); err != nil {
		return err
	}

	m.cpu.Wake()

	return nil
}