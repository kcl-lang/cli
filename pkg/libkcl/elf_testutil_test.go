// Copyright The KCL Authors. All rights reserved.

package libkcl

import (
	"debug/elf"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// writeTestLibrary writes a tiny ELF64 shared object that defines symbol.
// The file is only large enough for debug/elf to read the dynamic symbol
// table, so tests do not need a C toolchain.
func writeTestLibrary(t *testing.T, dir, symbol string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, LibraryFile)
	if err := os.WriteFile(path, testELF(symbol), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func testELF(symbol string) []byte {
	dynsym := make([]byte, 24*2)
	// Symbol 1 is the exported function. Symbol 0 is the required null entry.
	nameOff := 1 // dynstr starts with a NUL
	binary.LittleEndian.PutUint32(dynsym[24:], uint32(nameOff))
	dynsym[24+4] = byte(elf.STB_GLOBAL)<<4 | byte(elf.STT_FUNC)
	binary.LittleEndian.PutUint16(dynsym[24+6:], 1) // defined, not SHN_UNDEF

	dynstr := append([]byte{0}, append([]byte(symbol), 0)...)
	shstr := []byte("\x00.dynsym\x00.dynstr\x00.shstrtab\x00")

	const (
		ehSize = 64
		shSize = 64
	)
	dynsymOff := ehSize
	dynstrOff := dynsymOff + len(dynsym)
	shstrOff := dynstrOff + len(dynstr)
	shOff := shstrOff + len(shstr)

	buf := make([]byte, shOff+4*shSize)
	copy(buf[dynsymOff:], dynsym)
	copy(buf[dynstrOff:], dynstr)
	copy(buf[shstrOff:], shstr)

	// ELF header.
	copy(buf, []byte{0x7f, 'E', 'L', 'F', 2, 1, 1})
	binary.LittleEndian.PutUint16(buf[16:], uint16(elf.ET_DYN))
	binary.LittleEndian.PutUint16(buf[18:], uint16(elf.EM_X86_64))
	binary.LittleEndian.PutUint32(buf[20:], 1)
	binary.LittleEndian.PutUint64(buf[40:], uint64(shOff))
	binary.LittleEndian.PutUint16(buf[52:], ehSize)
	binary.LittleEndian.PutUint16(buf[58:], shSize)
	binary.LittleEndian.PutUint16(buf[60:], 4)
	binary.LittleEndian.PutUint16(buf[62:], 3)

	putSh := func(i int, name uint32, typ elf.SectionType, off, size uint64, link, info uint32, entsize uint64) {
		b := buf[shOff+i*shSize : shOff+(i+1)*shSize]
		binary.LittleEndian.PutUint32(b[0:], name)
		binary.LittleEndian.PutUint32(b[4:], uint32(typ))
		binary.LittleEndian.PutUint64(b[24:], off)
		binary.LittleEndian.PutUint64(b[32:], size)
		binary.LittleEndian.PutUint32(b[40:], link)
		binary.LittleEndian.PutUint32(b[44:], info)
		binary.LittleEndian.PutUint64(b[48:], 8)
		binary.LittleEndian.PutUint64(b[56:], entsize)
	}
	// shstrtab layout: NUL, ".dynsym", NUL, ".dynstr", NUL, ".shstrtab", NUL
	putSh(1, 1, elf.SHT_DYNSYM, uint64(dynsymOff), uint64(len(dynsym)), 2, 1, 24)
	putSh(2, 9, elf.SHT_STRTAB, uint64(dynstrOff), uint64(len(dynstr)), 0, 0, 0)
	putSh(3, 17, elf.SHT_STRTAB, uint64(shstrOff), uint64(len(shstr)), 0, 0, 0)
	return buf
}
