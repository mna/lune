package serializer

import (
	"encoding/binary"
	"fmt"
	"github.com/PuerkitoBio/lune/types"
	"io"
	"unsafe"
)

// TODO : Tests, and error variables much like io.EOF

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

func (h *gHeader) MajorVersion() byte {
	return h.Version >> 4
}

func (h *gHeader) MinorVersion() byte {
	return h.Version & 0x0F
}

func NewHeader() *gHeader {
	var i int32   // Force 4 bytes even on 64bit systems? Validate on 64bit Linux
	var ui uint64 // TODO : Force 8 bytes, uint gives only 4, inconsistent with Lua on 64-bit platforms
	var instr types.Instruction

	// Create a standard header based on the current architecture
	return &gHeader{
		Signature:  _HEADER,
		Version:    _VERSION,
		Format:     _FORMAT,
		Endianness: 1, // TODO : For now, force little-endian
		IntSz:      byte(unsafe.Sizeof(i)),
		SizeTSz:    byte(unsafe.Sizeof(ui)), // TODO : Is this consistent with what Lua gives on this platform?
		InstrSz:    byte(unsafe.Sizeof(instr)),
		NumberSz:   8, // TODO : Sizeof(the custom Number size)
		IntFlag:    0, // TODO : Support non-floating point compilation?
		Tail:       _TAIL,
	}
}

func readString(r io.Reader) string {
	var sz uint64
	var s string

	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		panic(err)
	}
	if sz > 0 {
		ch := make([]byte, sz)
		if err := binary.Read(r, binary.LittleEndian, ch); err != nil {
			panic(err)
		}
		// Remove 0x00
		s = string(ch[:len(ch)-1])
	}
	return s
}

func readDebug(r io.Reader, p *types.Prototype) {
	var n uint32
	var i uint32

	// Source file name
	p.Source = readString(r)

	// Line numbers
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		panic(err)
	}
	for i = 0; i < n; i++ {
		var li int32
		if err := binary.Read(r, binary.LittleEndian, &li); err != nil {
			panic(err)
		}
		p.LineInfo = append(p.LineInfo, li)
	}

	// Local variables
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		panic(err)
	}
	for i = 0; i < n; i++ {
		var lv types.LocVar
		lv.Name = readString(r)

		if err := binary.Read(r, binary.LittleEndian, &lv.Startpc); err != nil {
			panic(err)
		}
		if err := binary.Read(r, binary.LittleEndian, &lv.Endpc); err != nil {
			panic(err)
		}
		p.LocVars = append(p.LocVars, &lv)
	}

	// Upvalue names
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		panic(err)
	}
	for i = 0; i < n; i++ {
		p.Upvalues[i].Name = readString(r)
	}
}

func readUpvalues(r io.Reader, p *types.Prototype) {
	var n uint32
	var i uint32

	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		panic(err)
	}
	for i = 0; i < n; i++ {
		var ba [2]byte
		if err := binary.Read(r, binary.LittleEndian, &ba); err != nil {
			panic(err)
		}
		p.Upvalues = append(p.Upvalues, &types.Upvalue{"", ba[0], ba[1]})
	}
}

func readConstants(r io.Reader, p *types.Prototype) {
	var n uint32
	var i uint32

	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		panic(err)
	}

	for i = 0; i < n; i++ {
		// Read the constant's type, 1 byte
		var t byte
		if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
			panic(err)
		}
		switch types.ValType(t) {
		case types.TNIL:
			p.Ks = append(p.Ks, nil)
		case types.TBOOL:
			var v types.Value
			if err := binary.Read(r, binary.LittleEndian, &t); err != nil {
				panic(err)
			}
			if t == 0 {
				v = false
			} else if t == 1 {
				v = true
			} else {
				panic(fmt.Errorf("invalid value for boolean: %d", t))
			}
			p.Ks = append(p.Ks, v)
		case types.TNUMBER:
			var f float64
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				panic(err)
			}
			p.Ks = append(p.Ks, f)
		case types.TSTRING:
			p.Ks = append(p.Ks, readString(r))
		default:
			panic(fmt.Errorf("unexpected constant type: %d", t))
		}
	}
}

func readCode(r io.Reader, p *types.Prototype) {
	var n uint32
	var i uint32

	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		panic(err)
	}
	for i = 0; i < n; i++ {
		var instr types.Instruction
		if err := binary.Read(r, binary.LittleEndian, &instr); err != nil {
			panic(err)
		}
		p.Code = append(p.Code, instr)
	}
}

func readFunction(r io.Reader) *types.Prototype {
	var fm types.FuncMeta
	var p types.Prototype
	var n uint32
	var i uint32

	// Meta-data about the function
	if err := binary.Read(r, binary.LittleEndian, &fm); err != nil {
		panic(err)
	}
	p.Meta = &fm

	// Function's instructions
	readCode(r, &p)
	// Function's constants
	readConstants(r, &p)
	// Inner function's functions (prototypes)
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		panic(err)
	}
	for i = 0; i < n; i++ {
		p.Protos = append(p.Protos, readFunction(r))
	}

	// Upvalues
	readUpvalues(r, &p)
	// Debug
	readDebug(r, &p)

	return &p
}

func readHeader(r io.Reader) *gHeader {
	var h gHeader

	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		panic(err)
	}

	// Validate header
	stdH := NewHeader()

	// As a whole
	if h == *stdH {
		return &h
	} else if h.Signature != stdH.Signature {
		panic(fmt.Errorf("is not a precompiled chunk"))
	} else if h.Version != stdH.Version {
		panic(fmt.Errorf("version mismatch, got %d.%d, expected %d.%d", h.MajorVersion(), h.MinorVersion(), stdH.MajorVersion(), stdH.MinorVersion()))
	}

	panic(fmt.Errorf("incompatible"))
}

func Load(r io.Reader) (p *types.Prototype, err error) {
	// For simplicity's sake, to avoid multiple if err != nil, use a panic in the
	// various readXxxx functions, and catch here, since Load() returns as soon as an
	// error is detected (this is not a compiler).
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	// First up, the Header (12 bytes) + LUAC_TAIL to "catch conversion errors", as described in Lua
	readHeader(r)
	//fmt.Printf("Header: %+v\n", h)

	// Then, the function header (a prototype)
	p = readFunction(r)
	//fmt.Println(p)

	return
}
