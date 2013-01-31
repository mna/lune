package types

import (
	"bytes"
	"fmt"
	"strings"
)

/*
  Port of lopcodes.{c,h}

  Opcodes for the virtual machine
  See doc.go for copyright notice
*/

/*
  We assume that instructions are unsigned numbers.
  All instructions have an opcode in the first 6 bits.
  Instructions can have the following fields:
  `A' : 8 bits
  `B' : 9 bits
  `C' : 9 bits
  'Ax' : 26 bits ('A', 'B', and 'C' together)
  `Bx' : 18 bits (`B' and `C' together)
  `sBx' : signed Bx

  A signed argument is represented in excess K; that is, the number
  value is the unsigned value minus K. K is exactly the maximum value
  for that argument (so that -max is represented by 0, and +max is
  represented by 2*max), which is half the maximum for the corresponding
  unsigned argument.
*/
type Instruction uint32

func mask0(n, p uint) Instruction {
	return ^(mask1(n, p))
}

func mask1(n, p uint) Instruction {
	return ((^((^Instruction(0)) << n)) << p)
}

func getArg(i Instruction, pos, size uint) int {
	return int((i >> pos) & (mask1(size, 0)))
}

func (i Instruction) GetOpCode() OpCode {
	return OpCode((i >> posOp) & mask1(sizeOp, 0))
}

func (i Instruction) GetArgA() int {
	return getArg(i, posA, sizeA)
}

func getArgPossibleK(i Instruction, argNm byte, pos, size uint) int {
	var md OpArgMask

	arg := getArg(i, pos, size)
	op := i.GetOpCode()
	if argNm == 'B' {
		md = op.GetBMode()
	} else {
		md = op.GetCMode()
	}
	if md == OpArgK && isK(arg) {
		arg = indexK(arg)
	} else if md == OpArgN {
		panic(fmt.Sprintf("unexpected use of %s in operator %s", argNm, op))
	}
	return arg
}

func (i Instruction) GetArgB(getK bool) int {
	if getK {
		return getArgPossibleK(i, 'B', posB, sizeB)
	}
	return getArg(i, posB, sizeB)
}

func (i Instruction) GetArgC(getK bool) int {
	if getK {
		return getArgPossibleK(i, 'C', posC, sizeC)
	}
	return getArg(i, posC, sizeC)
}

func (i Instruction) GetArgBx(getK bool) int {
	if getK {
		return getArgPossibleK(i, 'B', posBx, sizeBx)
	}
	return getArg(i, posBx, sizeBx)
}

func (i Instruction) GetArgAx() int {
	return getArg(i, posAx, sizeAx)
}

func (i Instruction) GetArgsBx() int {
	return (i.GetArgBx(false) - MAXARG_sBx)
}

// test whether value is a constant
func isK(v int) bool {
	return (v & BITRK) != 0
}

// gets the index of the constant
func indexK(v int) int {
	return (v & (^BITRK))
}

func (i Instruction) String() string {
	var buf bytes.Buffer
	var a, b, c int
	var bc [2]bool

	op := i.GetOpCode()
	om := op.GetOpMode()

	buf.WriteString(fmt.Sprintf("%s%s", op, strings.Repeat(" ", 20-len(op.String()))))
	if om == MODE_iAx {
		a = i.GetArgAx()
	} else {
		a = i.GetArgA()

		bm := op.GetBMode()
		if bm != OpArgN {
			bc[0] = true
			switch om {
			case MODE_iABC:
				b = i.GetArgB(false)
				cm := op.GetCMode()
				if cm != OpArgN {
					bc[1] = true
					c = i.GetArgC(false)
					if cm == OpArgK && isK(c) {
						c = -1 - indexK(c)
					}
				}
			case MODE_iABx:
				b = i.GetArgBx(false)
			case MODE_iAsBx:
				b = i.GetArgsBx()
			}
			if bm == OpArgK && isK(b) {
				b = -1 - indexK(b)
			}
		}
	}
	buf.WriteString(fmt.Sprintf("%d", a))
	if bc[0] {
		buf.WriteString(fmt.Sprintf(" %d", b))
	}
	if bc[1] {
		buf.WriteString(fmt.Sprintf(" %d", c))
	}
	return buf.String()
}

// VM opcodes
type OpCode uint8

const (
	OP_MOVE OpCode = iota
	OP_LOADK
	OP_LOADKx
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETTABUP
	OP_GETTABLE
	OP_SETTABUP
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE
	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_POW
	OP_UNM
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TAILCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORCALL
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSURE
	OP_VARARG
	OP_EXTRAARG
	op_count
)

var opNames = [...]string{
	OP_MOVE:     "MOVE",
	OP_LOADK:    "LOADK",
	OP_LOADKx:   "LOADKX",
	OP_LOADBOOL: "LOADBOOL",
	OP_LOADNIL:  "LOADNIL",
	OP_GETUPVAL: "GETUPVAL",
	OP_GETTABUP: "GETTABUP",
	OP_GETTABLE: "GETTABLE",
	OP_SETTABUP: "SETTABUP",
	OP_SETUPVAL: "SETUPVAL",
	OP_SETTABLE: "SETTABLE",
	OP_NEWTABLE: "NEWTABLE",
	OP_SELF:     "SELF",
	OP_ADD:      "ADD",
	OP_SUB:      "SUB",
	OP_MUL:      "MUL",
	OP_DIV:      "DIV",
	OP_MOD:      "MOD",
	OP_POW:      "POW",
	OP_UNM:      "UNM",
	OP_NOT:      "NOT",
	OP_LEN:      "LEN",
	OP_CONCAT:   "CONCAT",
	OP_JMP:      "JMP",
	OP_EQ:       "EQ",
	OP_LT:       "LT",
	OP_LE:       "LE",
	OP_TEST:     "TEST",
	OP_TESTSET:  "TESTSET",
	OP_CALL:     "CALL",
	OP_TAILCALL: "TAILCALL",
	OP_RETURN:   "RETURN",
	OP_FORLOOP:  "FORLOOP",
	OP_FORPREP:  "FORPREP",
	OP_TFORCALL: "TFORCALL",
	OP_TFORLOOP: "TFORLOOP",
	OP_SETLIST:  "SETLIST",
	OP_CLOSURE:  "CLOSURE",
	OP_VARARG:   "VARARG",
	OP_EXTRAARG: "EXTRAARG",
}

func (o OpCode) String() string {
	return opNames[o]
}

/*
TODO :
#define testAMode(m)	(luaP_opmodes[m] & (1 << 6))
#define testTMode(m)	(luaP_opmodes[m] & (1 << 7))
*/

func (o OpCode) GetOpMode() OpMode {
	return OpMode(opMasks[o] & 3)
}

func (o OpCode) GetBMode() OpArgMask {
	return OpArgMask((opMasks[o] >> 4) & 3)
}

func (o OpCode) GetCMode() OpArgMask {
	return OpArgMask((opMasks[o] >> 2) & 3)
}

// Operator mode, defines how to access the other bits of the instruction.
type OpMode byte

const (
	MODE_iABC OpMode = iota
	MODE_iABx
	MODE_iAsBx
	MODE_iAx
)

/*
** masks for instruction properties. The format is:
** bits 0-1: op mode
** bits 2-3: C arg mode
** bits 4-5: B arg mode
** bit 6: instruction set register A
** bit 7: operator is a test (next instruction must be a jump)
 */
type OpArgMask byte

const (
	OpArgN OpArgMask = iota // Argument is not used
	OpArgU                  // Argument is used
	OpArgR                  // Argument is a register or a jump offset
	OpArgK                  // Argument is a constant or register/constant
)

// OpMask defines the behaviour of the instruction.
type OpMask byte

func createOpMask(tst, regA byte, bArgMode, cArgMode OpArgMask, om OpMode) OpMask {
	return OpMask((tst << 7) | (regA << 6) | (byte(bArgMode) << 4) | (byte(cArgMode) << 2) | byte(om))
}

var opMasks = [...]OpMask{
	OP_MOVE:     createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
	OP_LOADK:    createOpMask(0, 1, OpArgK, OpArgN, MODE_iABx),
	OP_LOADKx:   createOpMask(0, 1, OpArgN, OpArgN, MODE_iABx),
	OP_LOADBOOL: createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
	OP_LOADNIL:  createOpMask(0, 1, OpArgU, OpArgN, MODE_iABC),
	OP_GETUPVAL: createOpMask(0, 1, OpArgU, OpArgN, MODE_iABC),
	OP_GETTABUP: createOpMask(0, 1, OpArgU, OpArgK, MODE_iABC),
	OP_GETTABLE: createOpMask(0, 1, OpArgR, OpArgK, MODE_iABC),
	OP_SETTABUP: createOpMask(0, 0, OpArgK, OpArgK, MODE_iABC),
	OP_SETUPVAL: createOpMask(0, 0, OpArgU, OpArgN, MODE_iABC),
	OP_SETTABLE: createOpMask(0, 0, OpArgK, OpArgK, MODE_iABC),
	OP_NEWTABLE: createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
	OP_SELF:     createOpMask(0, 1, OpArgR, OpArgK, MODE_iABC),
	OP_ADD:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
	OP_SUB:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
	OP_MUL:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
	OP_DIV:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
	OP_MOD:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
	OP_POW:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
	OP_UNM:      createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
	OP_NOT:      createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
	OP_LEN:      createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
	OP_CONCAT:   createOpMask(0, 1, OpArgR, OpArgR, MODE_iABC),
	OP_JMP:      createOpMask(0, 0, OpArgR, OpArgN, MODE_iAsBx),
	OP_EQ:       createOpMask(1, 0, OpArgK, OpArgK, MODE_iABC),
	OP_LT:       createOpMask(1, 0, OpArgK, OpArgK, MODE_iABC),
	OP_LE:       createOpMask(1, 0, OpArgK, OpArgK, MODE_iABC),
	OP_TEST:     createOpMask(1, 0, OpArgN, OpArgU, MODE_iABC),
	OP_TESTSET:  createOpMask(1, 1, OpArgR, OpArgU, MODE_iABC),
	OP_CALL:     createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
	OP_TAILCALL: createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
	OP_RETURN:   createOpMask(0, 0, OpArgU, OpArgN, MODE_iABC),
	OP_FORLOOP:  createOpMask(0, 1, OpArgR, OpArgN, MODE_iAsBx),
	OP_FORPREP:  createOpMask(0, 1, OpArgR, OpArgN, MODE_iAsBx),
	OP_TFORCALL: createOpMask(0, 0, OpArgN, OpArgU, MODE_iABC),
	OP_TFORLOOP: createOpMask(0, 1, OpArgR, OpArgN, MODE_iAsBx),
	OP_SETLIST:  createOpMask(0, 0, OpArgU, OpArgU, MODE_iABC),
	OP_CLOSURE:  createOpMask(0, 1, OpArgU, OpArgN, MODE_iABx),
	OP_VARARG:   createOpMask(0, 1, OpArgU, OpArgN, MODE_iABC),
	OP_EXTRAARG: createOpMask(0, 0, OpArgU, OpArgU, MODE_iAx),
}
