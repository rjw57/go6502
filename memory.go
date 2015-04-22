/*
	ROM & RAM for go6502; 16-bit address, 8-bit data.
*/
package go6502

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

// Memory is a general interface for reading and writing bytes to and from
// 16-bit addresses.
type Memory interface {
	Shutdown()
	Read(uint16) byte
	Write(uint16, byte)
	Size() int
}

// A Rom provides read-only memory, with data generally pre-loaded from a file.
type Rom struct {
	name string
	size int // bytes
	data []byte
}

// Shutdown is part of the Memory interface, but takes no action for Rom.
func (r *Rom) Shutdown() {
}

// Read a byte from the given address.
func (rom *Rom) Read(a uint16) byte {
	return rom.data[a]
}

// Create a new ROM, loading the contents from a file.
// The size of the ROM is determined by the size of the file.
func RomFromFile(path string) (*Rom, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Rom{name: path, size: len(data), data: data}, nil
}

// Size of the Rom in bytes.
func (r *Rom) Size() int {
	return r.size
}

func (r *Rom) String() string {
	return fmt.Sprintf("ROM[%dk:%s:%s..%s]",
		r.Size()/1024,
		r.name,
		hex.EncodeToString(r.data[0:2]),
		hex.EncodeToString(r.data[len(r.data)-2:]))
}

// Rom meets the go6502.Memory interface, but Write is not supported, and will
// cause an error.
func (r *Rom) Write(_ uint16, _ byte) {
	panic(fmt.Sprintf("%v is read-only", r))
}

// Ram (32 KiB)
type Ram [0x8000]byte

// Shutdown is part of the Memory interface, but takes no action for Ram.
func (r *Ram) Shutdown() {
}

func (r *Ram) String() string {
	return "(RAM 32K)"
}

// Read a byte from a 16-bit address.
func (mem *Ram) Read(a uint16) byte {
	return mem[a]
}

// Write a byte to a 16-bit address.
func (mem *Ram) Write(a uint16, value byte) {
	mem[a] = value
}

// Size of the RAM in bytes.
func (mem *Ram) Size() int {
	return 0x8000 // 32K
}

// Dump writes the RAM contents to the specified file path.
func (mem *Ram) Dump(path string) {
	err := ioutil.WriteFile(path, mem[:], 0640)
	if err != nil {
		panic(err)
	}
}

// OffsetMemory wraps a Memory object, rewriting read/write addresses by the
// given offset. This makes it possible to mount memory into a larger address
// space at a given base address.
type OffsetMemory struct {
	Offset uint16
	Memory
}

// Read returns a byte from the underlying Memory after rewriting the address
// using the offset.
func (om OffsetMemory) Read(a uint16) byte {
	return om.Memory.Read(a - om.Offset)
}

func (om OffsetMemory) String() string {
	return fmt.Sprintf("OffsetMemory(%v)", om.Memory)
}

// Write stores a byte in the underlying Memory after rewriting the address
// using the offset.
func (om OffsetMemory) Write(a uint16, value byte) {
	om.Memory.Write(a-om.Offset, value)
}
