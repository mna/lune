package types

import (
	"strconv"
)

/*
  Mostly a port of lopcodes.{c,h} from Lua

  Opcodes for the virtual machine
  See doc.go for copyright notice
*/

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
	if int(o) >= len(opNames) {
		return strconv.Itoa(int(o))
	}
	return opNames[o]
}

type Args struct {
	Ax, Bx, Cx int
	A, B, C    *Value
	Bk, Ck     bool
}

type getArgsFunc func(*State, Instruction) Args

var opArgsFunc = [...]getArgsFunc{
	OP_MOVE:     getRARB,
	OP_LOADK:    getRAKBx,
	OP_LOADKx:   getRA,
	OP_LOADBOOL: getRABC,
	OP_LOADNIL:  getRARB,
	OP_GETUPVAL: getRAUB,
	OP_GETTABUP: getRAUBRKC,
	OP_GETTABLE: getRARBRKC,
	OP_SETTABUP: getUARKBRKC,
	OP_SETUPVAL: getRAUB,
	OP_SETTABLE: getRARKBRKC,
	OP_NEWTABLE: getRABC,
	OP_SELF:     getRARBRKC,
	OP_ADD:      getRARKBRKC,
	OP_SUB:      getRARKBRKC,
	OP_MUL:      getRARKBRKC,
	OP_DIV:      getRARKBRKC,
	OP_MOD:      getRARKBRKC,
	OP_POW:      getRARKBRKC,
	OP_UNM:      getRARB,
	OP_NOT:      getRARB,
	OP_LEN:      getRARB,
	OP_CONCAT:   getRARBRC,
	OP_JMP:      getAsBx,
	OP_EQ:       getARKBRKC,
	OP_LT:       getARKBRKC,
	OP_LE:       getARKBRKC,
	OP_TEST:     getRAC,
	OP_TESTSET:  getRARBC,
	OP_CALL:     getRABC,
	OP_TAILCALL: getRABC,
	OP_RETURN:   getRAB,
	OP_FORLOOP:  getRAsBx,
	OP_FORPREP:  getRAsBx,
	OP_TFORCALL: getRAC,
	OP_TFORLOOP: getRAsBx,
	OP_SETLIST:  getRABC,
	OP_CLOSURE:  getRABx,
	OP_VARARG:   getRAB,
	OP_EXTRAARG: getAx,
}

func getAx(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgAx()

	return ag
}

func getRA(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.A = &s.CI.Frame[ag.Ax]

	return ag
}

func getRAB(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.A = &s.CI.Frame[ag.Ax]
	ag.Bx, _ = i.GetArgB(false)

	return ag
}

func getRAC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.A = &s.CI.Frame[ag.Ax]
	ag.Cx, _ = i.GetArgC(false)

	return ag
}

func getRARB(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, _ = i.GetArgB(false)
	ag.A = &s.CI.Frame[ag.Ax]
	ag.B = &s.CI.Frame[ag.Bx]

	return ag
}

func getRARBC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, _ = i.GetArgB(false)
	ag.A = &s.CI.Frame[ag.Ax]
	ag.B = &s.CI.Frame[ag.Bx]
	ag.Cx, _ = i.GetArgC(false)

	return ag
}

func getRARBRC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, _ = i.GetArgB(false)
	ag.Cx, _ = i.GetArgC(false)
	ag.A = &s.CI.Frame[ag.Ax]
	ag.B = &s.CI.Frame[ag.Bx]
	ag.C = &s.CI.Frame[ag.Cx]

	return ag
}

func getRARBRKC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, _ = i.GetArgB(false)
	ag.Cx, ag.Ck = i.GetArgC(true)
	ag.A = &s.CI.Frame[ag.Ax]
	ag.B = &s.CI.Frame[ag.Bx]
	if ag.Ck {
		ag.C = &s.CI.Cl.P.Ks[ag.Cx]
	} else {
		ag.C = &s.CI.Frame[ag.Cx]
	}

	return ag
}

func getRABC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, _ = i.GetArgB(false)
	ag.Cx, _ = i.GetArgC(false)
	ag.A = &s.CI.Frame[ag.Ax]

	return ag
}

func getRAUB(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, _ = i.GetArgB(false)
	ag.A = &s.CI.Frame[ag.Ax]
	ag.B = &s.CI.Cl.UpVals[ag.Bx]

	return ag
}

func getRARKBRKC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, ag.Bk = i.GetArgB(true)
	ag.Cx, ag.Ck = i.GetArgC(true)
	ag.A = &s.CI.Frame[ag.Ax]
	if ag.Bk {
		ag.B = &s.CI.Cl.P.Ks[ag.Bx]
	} else {
		ag.B = &s.CI.Frame[ag.Bx]
	}
	if ag.Ck {
		ag.C = &s.CI.Cl.P.Ks[ag.Cx]
	} else {
		ag.C = &s.CI.Frame[ag.Cx]
	}

	return ag
}

func getARKBRKC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, ag.Bk = i.GetArgB(true)
	ag.Cx, ag.Ck = i.GetArgC(true)
	if ag.Bk {
		ag.B = &s.CI.Cl.P.Ks[ag.Bx]
	} else {
		ag.B = &s.CI.Frame[ag.Bx]
	}
	if ag.Ck {
		ag.C = &s.CI.Cl.P.Ks[ag.Cx]
	} else {
		ag.C = &s.CI.Frame[ag.Cx]
	}

	return ag
}

func getRAUBRKC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, _ = i.GetArgB(false)
	ag.Cx, ag.Ck = i.GetArgC(true)
	ag.A = &s.CI.Frame[ag.Ax]
	ag.B = &s.CI.Cl.UpVals[ag.Bx]
	if ag.Ck {
		ag.C = &s.CI.Cl.P.Ks[ag.Cx]
	} else {
		ag.C = &s.CI.Frame[ag.Cx]
	}

	return ag
}

func getUARKBRKC(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx, ag.Bk = i.GetArgB(true)
	ag.Cx, ag.Ck = i.GetArgC(true)
	ag.A = &s.CI.Cl.UpVals[ag.Ax]
	if ag.Bk {
		ag.B = &s.CI.Cl.P.Ks[ag.Bx]
	} else {
		ag.B = &s.CI.Frame[ag.Bx]
	}
	if ag.Ck {
		ag.C = &s.CI.Cl.P.Ks[ag.Cx]
	} else {
		ag.C = &s.CI.Frame[ag.Cx]
	}

	return ag
}

func getRAKBx(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.A = &s.CI.Frame[ag.Ax]

	// KBx implies that we read Bx as an int, and always look it up as a K
	ag.Bx, _ = i.GetArgBx(false)
	ag.Bk = true
	ag.B = &s.CI.Cl.P.Ks[ag.Bx]

	return ag
}

func getRABx(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.A = &s.CI.Frame[ag.Ax]
	ag.Bx, _ = i.GetArgBx(false)

	return ag
}

func getAsBx(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.Bx = i.GetArgsBx()

	return ag
}

func getRAsBx(s *State, i Instruction) Args {
	var ag Args

	ag.Ax = i.GetArgA()
	ag.A = &s.CI.Frame[ag.Ax]
	ag.Bx = i.GetArgsBx()

	return ag
}

// /*
// TODO :
// #define testAMode(m)	(luaP_opmodes[m] & (1 << 6))
// #define testTMode(m)	(luaP_opmodes[m] & (1 << 7))
// */

// func (o OpCode) GetOpMode() OpMode {
// 	return OpMode(opMasks[o] & 3)
// }

// func (o OpCode) GetAMode() bool {
// 	return (opMasks[o] & (1 << 6)) != 0
// }

// func (o OpCode) GetBMode() OpArgMask {
// 	return OpArgMask((opMasks[o] >> 4) & 3)
// }

// func (o OpCode) GetCMode() OpArgMask {
// 	return OpArgMask((opMasks[o] >> 2) & 3)
// }

// // Operator mode, defines how to access the other bits of the instruction.
// type OpMode byte

// const (
// 	MODE_iABC OpMode = iota
// 	MODE_iABx
// 	MODE_iAsBx
// 	MODE_iAx
// )

// /*
// ** masks for instruction properties. The format is:
// ** bits 0-1: op mode
// ** bits 2-3: C arg mode
// ** bits 4-5: B arg mode
// ** bit 6: instruction set register A
// ** bit 7: operator is a test (next instruction must be a jump)
//  */
// type OpArgMask byte

// const (
// 	OpArgN OpArgMask = iota // Argument is not used
// 	OpArgU                  // Argument is used
// 	OpArgR                  // Argument is a register or a jump offset
// 	OpArgK                  // Argument is a constant or register/constant
// )

// // OpMask defines the behaviour of the instruction.
// type OpMask byte

// func createOpMask(tst, regA byte, bArgMode, cArgMode OpArgMask, om OpMode) OpMask {
// 	return OpMask((tst << 7) | (regA << 6) | (byte(bArgMode) << 4) | (byte(cArgMode) << 2) | byte(om))
// }

// var opMasks = [...]OpMask{
// 	OP_MOVE:     createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
// 	OP_LOADK:    createOpMask(0, 1, OpArgK, OpArgN, MODE_iABx),
// 	OP_LOADKx:   createOpMask(0, 1, OpArgN, OpArgN, MODE_iABx),
// 	OP_LOADBOOL: createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
// 	OP_LOADNIL:  createOpMask(0, 1, OpArgU, OpArgN, MODE_iABC),
// 	OP_GETUPVAL: createOpMask(0, 1, OpArgU, OpArgN, MODE_iABC),
// 	OP_GETTABUP: createOpMask(0, 1, OpArgU, OpArgK, MODE_iABC),
// 	OP_GETTABLE: createOpMask(0, 1, OpArgR, OpArgK, MODE_iABC),
// 	OP_SETTABUP: createOpMask(0, 0, OpArgK, OpArgK, MODE_iABC),
// 	OP_SETUPVAL: createOpMask(0, 0, OpArgU, OpArgN, MODE_iABC),
// 	OP_SETTABLE: createOpMask(0, 0, OpArgK, OpArgK, MODE_iABC),
// 	OP_NEWTABLE: createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
// 	OP_SELF:     createOpMask(0, 1, OpArgR, OpArgK, MODE_iABC),
// 	OP_ADD:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
// 	OP_SUB:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
// 	OP_MUL:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
// 	OP_DIV:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
// 	OP_MOD:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
// 	OP_POW:      createOpMask(0, 1, OpArgK, OpArgK, MODE_iABC),
// 	OP_UNM:      createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
// 	OP_NOT:      createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
// 	OP_LEN:      createOpMask(0, 1, OpArgR, OpArgN, MODE_iABC),
// 	OP_CONCAT:   createOpMask(0, 1, OpArgR, OpArgR, MODE_iABC),
// 	OP_JMP:      createOpMask(0, 0, OpArgR, OpArgN, MODE_iAsBx),
// 	OP_EQ:       createOpMask(1, 0, OpArgK, OpArgK, MODE_iABC),
// 	OP_LT:       createOpMask(1, 0, OpArgK, OpArgK, MODE_iABC),
// 	OP_LE:       createOpMask(1, 0, OpArgK, OpArgK, MODE_iABC),
// 	OP_TEST:     createOpMask(1, 0, OpArgN, OpArgU, MODE_iABC),
// 	OP_TESTSET:  createOpMask(1, 1, OpArgR, OpArgU, MODE_iABC),
// 	OP_CALL:     createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
// 	OP_TAILCALL: createOpMask(0, 1, OpArgU, OpArgU, MODE_iABC),
// 	OP_RETURN:   createOpMask(0, 0, OpArgU, OpArgN, MODE_iABC),
// 	OP_FORLOOP:  createOpMask(0, 1, OpArgR, OpArgN, MODE_iAsBx),
// 	OP_FORPREP:  createOpMask(0, 1, OpArgR, OpArgN, MODE_iAsBx),
// 	OP_TFORCALL: createOpMask(0, 0, OpArgN, OpArgU, MODE_iABC),
// 	OP_TFORLOOP: createOpMask(0, 1, OpArgR, OpArgN, MODE_iAsBx),
// 	OP_SETLIST:  createOpMask(0, 0, OpArgU, OpArgU, MODE_iABC),
// 	OP_CLOSURE:  createOpMask(0, 1, OpArgU, OpArgN, MODE_iABx),
// 	OP_VARARG:   createOpMask(0, 1, OpArgU, OpArgN, MODE_iABC),
// 	OP_EXTRAARG: createOpMask(0, 0, OpArgU, OpArgU, MODE_iAx),
// }
