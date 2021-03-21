package core

import (
	"fmt"
)

type memory struct {
	data []Word
}

type MemoryOutOfBounds struct {
	addr Address
}

func (e *MemoryOutOfBounds) Error() string {
	return fmt.Sprintf("Address %d is out of bounds", e.addr)
}

func NewMemory(size uint) *memory {
	mem := new(memory)
	mem.data = make([]Word, size)
	return mem
}

func (m *memory) Load(offset Address, data Words) error {
	for addr, w := range data {
		if err := m.Write(offset + Address(addr), w); err != nil {
			return err
		}
	}

	return nil
}

func (m *memory) Read(addr Address) (error, Word) {
	if err := m.assertAccess(addr); err != nil {
		return err, 0
	}

	return nil, m.data[addr]
}

func (m *memory) Write(addr Address, data Word) error {
	if err := m.assertAccess(addr); err != nil {
		return err
	}

	m.data[addr] = data

	return nil
}

func (m *memory) assertAccess(addr Address) error {
	if addr >= Address(len(m.data)) {
		return &MemoryOutOfBounds{addr}
	}

	return nil
}