package serializer

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

const (
	LUNE_MAJOR_VERSION      = 5
	LUNE_MINOR_VERSION      = 2
	_VERSION           byte = LUNE_MAJOR_VERSION*16 + LUNE_MINOR_VERSION
	_FORMAT            byte = 0
	_HEADER_SZ              = 4
	_TAIL_SZ                = 6
)

var (
	_HEADER = [...]byte{0x1B, 0x4C, 0x75, 0x61}
	_TAIL   = [...]byte{0x19, 0x93, '\r', '\n', 0x1a, '\n'}
)

type gHeader struct {
	Signature  [_HEADER_SZ]byte
	Version    byte
	Format     byte
	Endianness byte
	IntSz      byte
	SizeTSz    byte
	InstrSz    byte
	NumberSz   byte
	IntFlag    byte
	Tail       [_TAIL_SZ]byte
}

func (h gHeader) MajorVersion() byte {
	return h.Version >> 4
}

func (h gHeader) MinorVersion() byte {
	return h.Version & 0x0F
}

func NewHeader() *gHeader {
	var i int
	var ui uint64 // Force 8 bytes, uint gives only 4, inconsistent with Lua on 64-bit platforms

	// Create a standard header based on the current architecture
	return &gHeader{
		Signature:  _HEADER,
		Version:    _VERSION,
		Format:     _FORMAT,
		Endianness: 1, // TODO : For now, force little-endian
		IntSz:      byte(unsafe.Sizeof(i)),
		SizeTSz:    byte(unsafe.Sizeof(ui)), // TODO : Is this consistent with what Lua gives on this platform?
		InstrSz:    4,                       // TODO : Sizeof(Instruction) in Lua
		NumberSz:   8,                       // TODO : Sizeof(the custom Number size)
		IntFlag:    0,                       // TODO : Support non-floating point compilation?
		Tail:       _TAIL,
	}
}

func readString(r io.Reader) (string, error) {
	var sz uint64
	var s string

	fmt.Println("Sizeof uint64 ", unsafe.Sizeof(sz))
	err := binary.Read(r, binary.LittleEndian, &sz)
	if err != nil {
		return "", err
	}
	if sz > 0 {
		fmt.Println("sz= ", sz)
		ch := make([]byte, sz)
		err = binary.Read(r, binary.LittleEndian, ch)
		if err != nil {
			return "", err
		}
		s = string(ch)
	}
	return s, nil
}

func readHeader(r io.Reader) (*gHeader, error) {
	var h gHeader

	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	// Validate header
	stdH := NewHeader()
	fmt.Printf("h: %v\n", h)
	fmt.Printf("stdH: %v\n", *stdH)
	if h != *stdH {
		return nil, fmt.Errorf("invalid header")
	}
	return &h, nil
}

func Load(r io.Reader) error {
	// First up, the Header (12 bytes)
	h, err := readHeader(r)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", h)

	// Then, the LUAC_TAIL to "catch conversion errors", as described in Lua

	s, err := readString(r)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)

	return nil
}
