package types

import (
	"fmt"
)

/*
  Port of lopcodes.h

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
type OpMode int

const (
	MODE_iABC OpMode = iota
	MODE_iABx
	MODE_iAsBx
	MODE_iAx
)

type Instruction int32

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

func (i Instruction) GetArgB() int {
	return getArg(i, posB, sizeB)
}

func (i Instruction) GetArgC() int {
	return getArg(i, posC, sizeC)
}

func (i Instruction) GetArgBx() int {
	return getArg(i, posBx, sizeBx)
}

func (i Instruction) GetArgAx() int {
	return getArg(i, posAx, sizeAx)
}

func (i Instruction) GetArgsBx() int {
	return (i.GetArgBx() - MAXARG_sBx)
}

// TODO : Needs more work to translate K notation, see lvm.c

func (i Instruction) String() string {
	// TODO : Switch depending on opmode
	return fmt.Sprintf("%s a=%d, b=%d, c=%d, ax=%d, bx=%d, sbx=%d", i.GetOpCode(), i.GetArgA(), i.GetArgB(), i.GetArgC(), i.GetArgAx(), i.GetArgBx(), i.GetArgsBx())
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
