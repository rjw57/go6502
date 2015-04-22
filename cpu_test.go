package go6502

import (
	"fmt"
	"testing"
)

func createCpu() *Cpu {
	ram := &Ram{}
	addressBus, _ := CreateBus()
	addressBus.Attach(ram, "ram", 0x8000) // upper 32K
	cpu := &Cpu{Bus: addressBus}
	cpu.Reset()
	return cpu
}

func TestBitInstruction(t *testing.T) {
	cpu := createCpu()
	cpu.Bus.Write(0x8000, 0xAA)

	instruction := Instruction{OpType: optypes[0x2C], Op16: 0x8000}
	expectedName := "BIT absolute $8000"
	actualName := instruction.String()
	if actualName != expectedName {
		t.Error(fmt.Sprintf("expected %s, got %s\n", expectedName, actualName))
	}

	cpu.BIT(instruction)

	expectedStatus := "n-_b-iz-"
	actualStatus := cpu.statusString()
	if actualStatus != expectedStatus {
		t.Error(fmt.Sprintf("SR expected %s got %s\n", expectedStatus, actualStatus))
	}
}
