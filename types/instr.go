package types

import (
	"bytes"
	"fmt"
	"strings"
)

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

func getArgPossibleK(i Instruction, argNm byte, pos, size uint) (int, bool) {
	var md OpArgMask
	var k bool

	arg := getArg(i, pos, size)
	op := i.GetOpCode()
	if argNm == 'B' {
		md = op.GetBMode()
	} else {
		md = op.GetCMode()
	}
	if md == OpArgK && isK(arg) {
		k = true
		arg = indexK(arg)
	} else if md == OpArgN {
		panic(fmt.Sprintf("unexpected use of %s in operator %s", argNm, op))
	}
	return arg, k
}

func (i Instruction) GetArgB(getK bool) (int, bool) {
	if getK {
		return getArgPossibleK(i, 'B', posB, sizeB)
	}
	return getArg(i, posB, sizeB), false
}

func (i Instruction) GetArgC(getK bool) (int, bool) {
	if getK {
		return getArgPossibleK(i, 'C', posC, sizeC)
	}
	return getArg(i, posC, sizeC), false
}

func (i Instruction) GetArgBx(getK bool) (int, bool) {
	if getK {
		return getArgPossibleK(i, 'B', posBx, sizeBx)
	}
	return getArg(i, posBx, sizeBx), false
}

func (i Instruction) GetArgAx() int {
	return getArg(i, posAx, sizeAx)
}

func (i Instruction) GetArgsBx() int {
	bx, _ := i.GetArgBx(false)
	return (bx - MAXARG_sBx)
}

// test whether value is a constant
func isK(v int) bool {
	return (v & BITRK) != 0
}

// gets the index of the constant
func indexK(v int) int {
	return (v & (^BITRK))
}

func (i Instruction) GetArgs(s *State) (a, b, c *Value) {
	var ax, bx, cx int

	op := i.GetOpCode()
	om := op.GetOpMode()

	if om == MODE_iAx {
		ax = i.GetArgAx()
	} else {
		ax = i.GetArgA()
		bm := op.GetBMode()
		if bm != OpArgN {
			switch om {
			case MODE_iABC:
				bx, _ = i.GetArgB(true)
			case MODE_iABx:
				bx, _ = i.GetArgBx(true)
			case MODE_iAsBx:
				bx = i.GetArgsBx()
			}

			switch bm {
			case OpArgK:
				b = &s.CI.Cl.P.Ks[bx]
			case OpArgR:
				b = &s.Frame[bx]
			case OpArgU:
				b = &s.CI.Cl.UpVals[bx]
			}
		}
		if om == MODE_iABC {
			cm := op.GetCMode()
			cx, _ = i.GetArgC(true)
			switch cm {
			case OpArgK:
				c = &s.CI.Cl.P.Ks[cx]
			case OpArgR:
				c = &s.Frame[cx]
			case OpArgU:
				c = &s.CI.Cl.UpVals[cx]
			}
		}
	}

	if op.GetAMode() {
		// Register
		a = &s.Frame[ax]
	} else {
		// Upvalue
		a = &s.CI.Cl.UpVals[ax]
	}
	return
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
				b, _ = i.GetArgB(false)
				cm := op.GetCMode()
				if cm != OpArgN {
					bc[1] = true
					c, _ = i.GetArgC(false)
					if cm == OpArgK && isK(c) {
						c = -1 - indexK(c)
					}
				}
			case MODE_iABx:
				b, _ = i.GetArgBx(false)
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
