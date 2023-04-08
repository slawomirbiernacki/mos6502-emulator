package memory

type MemoryMapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
}

type DummyMemoryMapper struct {
	Mem [65536]byte
}

func (m *DummyMemoryMapper) Read(address uint16) byte {
	return m.Mem[address]
}

func (m *DummyMemoryMapper) Write(address uint16, value byte) {
	m.Mem[address] = value
}
