package serializer

import (
	"encoding/binary"
	"fmt"
	"io"
)

var (
	HeaderSignature = [...]byte{0x1B, 0x4C, 0x75, 0x61}
)

type gHeader struct {
	Signature  [4]byte
	Version    byte
	Format     byte
	Endianness byte
	IntSz      byte
	SizeTSz    byte
	InstrSz    byte
	NumberSz   byte
	IntFlag    byte
}

func (h gHeader) MajorVersion() byte {
	return h.Version >> 4
}

func (h gHeader) MinorVersion() byte {
	return h.Version & 0x0F
}

func readHeader(r io.Reader) (*gHeader, error) {
	var h gHeader

	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	// Validate signature
	if h.Signature != HeaderSignature {
		return nil, fmt.Errorf("invalid signature")
	}
	return &h, nil
}

func Load(r io.Reader) error {
	_, err := readHeader(r)
	if err != nil {
		return err
	}
	return nil
}
