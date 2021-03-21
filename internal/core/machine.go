package core

type Machine struct {
	memory *memory
	cpu *cpu
}

func NewMachine(memSize uint, cpuFlags map[string]int) *Machine {
	m := NewMemory(memSize)

	return &Machine{memory: m, cpu: NewCpu(m, cpuFlags)}
}

func (m *Machine) Boot(program Words) error {
	if err := m.memory.Load(0, program); err != nil {
		return err
	}

	m.cpu.Wake()

	return nil
}