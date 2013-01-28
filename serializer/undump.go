package serializer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/PuerkitoBio/lune/types"
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

type prototype struct {
	meta     *funcMeta
	code     []types.Instruction
	ks       []types.Value
	protos   []*prototype
	upvalues []*upvalue

	// Debug info, unavailable in release build
	source   string
	lineInfo []int32
	locVars  []*locVar
}

func (p *prototype) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%+v\n", p.meta))
	buf.WriteString(fmt.Sprintln("Instructions (", len(p.code), ") :"))
	for _, c := range p.code {
		buf.WriteString(fmt.Sprintln(c))
	}
	buf.WriteString(fmt.Sprintln("Constants (", len(p.ks), ") :"))
	buf.WriteString(fmt.Sprintln(p.ks))
	buf.WriteString(fmt.Sprintln("Functions (", len(p.protos), ") :"))
	for _, f := range p.protos {
		buf.WriteString(fmt.Sprintln(f))
	}
	buf.WriteString(fmt.Sprintln("Upvalues (", len(p.upvalues), ") :"))
	for _, u := range p.upvalues {
		buf.WriteString(fmt.Sprintf("%+v\n", u))
	}
	buf.WriteString("\nDebug information:\n\n")
	buf.WriteString("Source: " + p.source + "\n")
	buf.WriteString(fmt.Sprintln("Line info (", len(p.lineInfo), ") :"))
	buf.WriteString(fmt.Sprintln(p.lineInfo))
	buf.WriteString(fmt.Sprintln("Local variables (", len(p.locVars), ") :"))
	for _, lv := range p.locVars {
		buf.WriteString(fmt.Sprintf("%+v\n", lv))
	}

	return buf.String()
}

type funcMeta struct {
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVarArg        byte
	MaxStackSize    byte
}

type upvalue struct {
	name    string
	instack byte
	idx     byte
}

type locVar struct {
	name    string
	startpc int
	endpc   int
}

func readString(r io.Reader) (string, error) {
	var sz uint64
	var s string

	err := binary.Read(r, binary.LittleEndian, &sz)
	if err != nil {
		return "", err
	}
	if sz > 0 {
		ch := make([]byte, sz)
		err = binary.Read(r, binary.LittleEndian, ch)
		if err != nil {
			return "", err
		}
		// Remove 0x00
		s = string(ch[:len(ch)-1])
	}
	return s, nil
}

func readDebug(r io.Reader, p *prototype) error {
	var n uint32
	var i uint32
	var err error

	// Source file name
	p.source, err = readString(r)
	if err != nil {
		return err
	}

	// Line numbers
	err = binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return err
	}
	for i = 0; i < n; i++ {
		var li int32
		err = binary.Read(r, binary.LittleEndian, &li)
		if err != nil {
			return err
		}
		p.lineInfo = append(p.lineInfo, li)
	}

	// Local variables
	err = binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return err
	}
	for i = 0; i < n; i++ {
		var lv locVar
		lv.name, err = readString(r)
		if err != nil {
			return err
		}
		err = binary.Read(r, binary.LittleEndian, &lv.startpc)
		if err != nil {
			return err
		}
		err = binary.Read(r, binary.LittleEndian, &lv.endpc)
		if err != nil {
			return err
		}
		p.locVars = append(p.locVars, &lv)
	}

	// Upvalue names
	err = binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return err
	}
	for i = 0; i < n; i++ {
		p.upvalues[i].name, err = readString(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func readUpvalues(r io.Reader, p *prototype) error {
	var n uint32
	var i uint32

	err := binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return err
	}
	for i = 0; i < n; i++ {
		var ba [2]byte
		err = binary.Read(r, binary.LittleEndian, &ba)
		if err != nil {
			return err
		}
		p.upvalues = append(p.upvalues, &upvalue{"", ba[0], ba[1]})
	}

	return nil
}

func readConstants(r io.Reader, p *prototype) error {
	var n uint32
	var i uint32

	err := binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return err
	}

	for i = 0; i < n; i++ {
		// Read the constant's type, 1 byte
		var t byte
		err = binary.Read(r, binary.LittleEndian, &t)
		if err != nil {
			return err
		}
		switch types.ValType(t) {
		case types.TNIL:
			var v types.Value = nil
			p.ks = append(p.ks, v)
		case types.TBOOL:
			var v types.Value
			err = binary.Read(r, binary.LittleEndian, &t)
			if err != nil {
				return err
			}
			if t == 0 {
				v = false
			} else if t == 1 {
				v = true
			} else {
				return fmt.Errorf("invalid value for boolean: %d", t)
			}
			p.ks = append(p.ks, v)
		case types.TNUMBER:
			// TODO : A number is a double in Lua, will a read in a float64 work?
			var f float64
			var v types.Value
			err = binary.Read(r, binary.LittleEndian, &f)
			if err != nil {
				return err
			}
			v = f
			p.ks = append(p.ks, v)
		case types.TSTRING:
			var v types.Value
			v, err = readString(r)
			if err != nil {
				return err
			}
			p.ks = append(p.ks, v)
		default:
			return fmt.Errorf("unexpected constant type: %d", t)
		}
	}
	return nil
}

func readCode(r io.Reader, p *prototype) error {
	var n uint32
	var i uint32

	err := binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return err
	}
	for i = 0; i < n; i++ {
		var instr types.Instruction
		err = binary.Read(r, binary.LittleEndian, &instr)
		if err != nil {
			return err
		}
		p.code = append(p.code, instr)
	}
	return nil
}

func readFunction(r io.Reader) (*prototype, error) {
	var fm funcMeta
	var p prototype
	var n uint32
	var i uint32

	// Meta-data about the function
	err := binary.Read(r, binary.LittleEndian, &fm)
	if err != nil {
		return nil, err
	}
	p.meta = &fm

	// Function's instructions
	err = readCode(r, &p)
	if err != nil {
		return nil, err
	}

	// Function's constants
	err = readConstants(r, &p)
	if err != nil {
		return nil, err
	}

	// Inner function's functions (prototypes)
	err = binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return nil, err
	}
	for i = 0; i < n; i++ {
		var subP *prototype
		subP, err = readFunction(r)
		if err != nil {
			return nil, err
		}
		p.protos = append(p.protos, subP)
	}

	// Upvalues
	err = readUpvalues(r, &p)
	if err != nil {
		return nil, err
	}

	// Debug
	err = readDebug(r, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func readHeader(r io.Reader) (*gHeader, error) {
	var h gHeader

	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	// Validate header
	stdH := NewHeader()

	// As a whole
	if h == *stdH {
		return &h, nil
	} else if h.Signature != stdH.Signature {
		return nil, fmt.Errorf("is not a precompiled chunk")
	} else if h.Version != stdH.Version {
		return nil, fmt.Errorf("version mismatch, got %d.%d, expected %d.%d", h.MajorVersion(), h.MinorVersion(), stdH.MajorVersion(), stdH.MinorVersion())
	}

	return nil, fmt.Errorf("incompatible")
}

func Load(r io.Reader) error {
	// First up, the Header (12 bytes) + LUAC_TAIL to "catch conversion errors", as described in Lua
	h, err := readHeader(r)
	if err != nil {
		return err
	}
	fmt.Printf("Header: %+v\n", h)

	// Then, the function header (a prototype)
	p, err := readFunction(r)
	if err != nil {
		return err
	}
	fmt.Println(p)

	return nil
}
