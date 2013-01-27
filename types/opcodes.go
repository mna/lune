package types

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

// Size of opcode arguments
const (
	sizeC  = 9
	sizeB  = 9
	sizeBx = sizeC + sizeB
	sizeA  = 8
	sizeAx = sizeC + sizeB + sizeA
	sizeOp = 6
)

// Position of opcode arguments
const (
	posOp = 0
	posA  = posOp + sizeOp
	posC  = posA + sizeA
	posB  = posC + sizeC
	posBx = posC
	posAx = posA
)

func GetOpCode(i Instruction) OpCode {
	var o OpCode = OpCode((i >> posOp) & (^(0 << sizeOp)))
	//#define GET_OPCODE(i)  (cast(OpCode, ((i)>>POS_OP) & MASK1(SIZE_OP,0)))
	return o
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
