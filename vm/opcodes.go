package vm

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
type opmode int

const (
	iABC opmode = iota
	iABx
	iAsBx
	iAx
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

// VM opcodes
type opcode int

const (
	opMove opcode = iota
	opLoadK
	opLoadKx
	opLoadBool
	opLoadNil
	opGetUpval
	opGetTabUp
	opGetTable
	opSetTabUp
	opSetUpval
	opSetTable
	opNewTable
	opSelf
	opAdd
	opSub
	opMul
	opDiv
	opMod
	opPow
	opUnm
	opNot
	opLen
	opConcat
	opJmp
	opEq
	opLt
	opLe
	opTest
	opTestSet
	opCall
	opTailCall
	opReturn
	opForLoop
	opForPrep
	opTForCall
	opTForLoop
	opSetList
	opClosure
	opVarArg
	opExtraArg
	op_count
)
