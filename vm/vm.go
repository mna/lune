package vm

import (
	"fmt"
	"github.com/PuerkitoBio/lune/types"
)

var (
	_BINOPS = [...]byte{
		types.OP_ADD: '+',
		types.OP_SUB: '-',
		types.OP_MUL: '*',
		types.OP_DIV: '/',
		types.OP_MOD: '%',
		types.OP_POW: '^',
	}

	_CMPOPS = [...]func(types.Value, types.Value) bool{
		types.OP_EQ: areEqual,
		types.OP_LT: isLessThan,
		types.OP_LE: isLessEqual,
	}
)

func doJump(s *types.State, i types.Instruction, e int) (ax, bx int) {
	ax = i.GetArgA()
	if ax != 0 {
		// TODO : Close upvalues - See dojump in lvm.c.
	}
	bx = i.GetArgsBx()
	s.CI.PC += bx + e
	return
}

func callGoFunc(s *types.State, f types.GoFunc, base, nRets int) {
	var in []types.Value
	for i := base; i < s.Top; i++ {
		in = append(in, s.Stack[i])
	}
	out := f(in)
	// Out values replace the stack values starting at the Go Func index (base - 1)
	// nRets values are expected, stop at this count, and fill with nils if necessary
	s.Top = base - 1
	for i := 0; i < nRets; i++ {
		if i < len(out) {
			s.Stack[s.Top] = out[i]
			s.Top++
		} else {
			s.Stack[s.Top] = nil
			s.Top++
		}
	}
}

func Execute(s *types.State) {
	var a, b, c *types.Value

	// Start with entry point (position 0)
	s.NewCallInfo(s.Stack[0].(*types.Closure), 0, 0)

newFrame:
	for {
		i := s.CI.Cl.P.Code[s.CI.PC]
		op := i.GetOpCode()
		s.CI.PC++
		s.DumpStack()
		// TODO : Will fail big time here, when some opcodes are executed. OpMasks
		// just doesn't do what I think it does.
		// TODO : Make OpMasks smarter for my needs, and return ax, bx, cx too?
		a, b, c = i.GetArgs(s)

		switch op {
		case types.OP_MOVE:
			// A B | R(A) := R(B)
			*a = *b
			fmt.Printf("%s\tR(A)=%v R(B)=%v\n", op, *a, *b)

		case types.OP_LOADK:
			// A Bx | R(A) := Kst(Bx)
			*a = *b
			fmt.Printf("%s\tR(A)=%v Kst(Bx)=%v\n", op, *a, *b)

		case types.OP_LOADKx:
			// A | R(A) := Kst(extra arg)
			// Special instruction: must always be followed by OP_EXTRAARG
			if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_EXTRAARG {
				panic(fmt.Sprintf("%s: expected OP_EXTRAARG as next instruction, found %s", op, i2.GetOpCode()))
			} else {
				s.CI.PC++
				ax := i2.GetArgAx()
				*a = s.CI.Cl.P.Ks[ax]
				fmt.Printf("%s\tR(A)=%v EXTRAARG=%v\n", op, *a, ax)
			}

		case types.OP_LOADBOOL:
			// A B C | R(A) := (Bool)B; if (C) PC++
			bx, _ := i.GetArgB(false)
			bb := bx != 0
			*a = bb

			// Skip next instruction if C is true
			cx, _ := i.GetArgC(false)
			cb := cx != 0
			if cb {
				s.CI.PC++
			}
			fmt.Printf("%s\tR(A)=%v B=%v C=%v\n", op, *a, bb, cb)

		case types.OP_LOADNIL:
			// A B | R(A) := ... := R(B) := nil
			ax := i.GetArgA()
			bx, _ := i.GetArgB(false)
			for j := 0; j <= bx; j++ {
				s.CI.Frame[ax+j] = nil
			}
			fmt.Printf("%s\tA=%v B=%v\n", op, ax, bx)

		case types.OP_GETUPVAL:
			// A B | R(A) := UpValue[B]
			*a = *b
			fmt.Printf("%s\tR(A)=%v U(B)=%v\n", op, *a, *b)

		case types.OP_GETTABUP:
			// A B C | R(A) := UpValue[B][RK(C)]
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s\tR(A)=%v U(B)=%v RK(C)=%v\n", op, *a, t, *c)

		case types.OP_GETTABLE:
			// A B C | R(A) := R(B)[RK(C)]
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s\tR(A)=%v R(B)=%v RK(C)=%v\n", op, *a, t, *c)

		case types.OP_SETTABUP:
			// A B C | UpValue[A][RK(B)] := RK(C) 
			t := (*a).(types.Table)
			t.Set(*b, *c)
			fmt.Printf("%s\tU(A)=%v RK(B)=%v RK(C)=%v\n", op, t, *b, *c)

		case types.OP_SETUPVAL:
			// A B | UpValue[B] := R(A)
			*b = *a
			fmt.Printf("%s\tR(A)=%v U(B)=%v\n", op, *a, *b)

		case types.OP_SETTABLE:
			// A B C | R(A)[RK(B)] := RK(C)
			t := (*a).(types.Table)
			t.Set(*b, *c)
			fmt.Printf("%s\tR(A)=%v RK(B)=%v RK(C)=%v\n", op, t, *b, *c)

		case types.OP_NEWTABLE:
			// A B C | R(A) := {} (size = B,C)
			t := types.NewTable()
			bx, _ := i.GetArgB(false)
			cx, _ := i.GetArgC(false)
			// TODO : Encoded array and hash sizes (B and C) are ignored at the moment
			*a = t
			fmt.Printf("%s\tR(A)=%v B=%v C=%v\n", op, t, bx, cx)

		case types.OP_SELF:
			// A B C | R(A+1) := R(B); R(A) := R(B)[RK(C)]
			ax := i.GetArgA()
			s.CI.Frame[ax+1] = *b
			t := (*b).(types.Table)
			s.CI.Frame[ax] = t.Get(*c)
			fmt.Printf("%s\tA=%v R(B)=%v RK(C)=%v\n", op, ax, t, *c)

		case types.OP_ADD, types.OP_SUB, types.OP_MUL, types.OP_DIV,
			types.OP_MOD, types.OP_POW:
			// A B C | R(A) := RK(B) + RK(C)
			// A B C | R(A) := RK(B) - RK(C)
			// A B C | R(A) := RK(B) * RK(C)
			// A B C | R(A) := RK(B) ÷ RK(C)
			// A B C | R(A) := RK(B) % RK(C)
			// A B C | R(A) := RK(B) ^ RK(C)
			*a = coerceAndComputeBinaryOp(_BINOPS[op], *b, *c)
			fmt.Printf("%s\tR(A)=%v RK(B)=%v RK(C)=%v\n", op, *a, *b, *c)

		case types.OP_UNM:
			// A B | R(A) := -R(B)
			*a = coerceAndComputeUnaryOp('-', *b)
			fmt.Printf("%s\tR(A)=%v R(B)=%v\n", op, *a, *b)

		case types.OP_NOT:
			// A B | R(A) := not R(B)
			*a = isFalse(*b)
			fmt.Printf("%s\tR(A)=%v R(B)=%v\n", op, *a, *b)

		case types.OP_LEN:
			// A B | R(A) := length of R(B)
			*a = computeLength(*b)
			fmt.Printf("%s\tR(A)=%v R(B)=%v\n", op, *a, *b)

		case types.OP_CONCAT:
			// A B C | R(A) := R(B).. ... ..R(C)
			bx, _ := i.GetArgB(false)
			cx, _ := i.GetArgC(false)
			src := s.CI.Frame[bx : cx+1]
			*a = coerceAndConcatenate(src)
			fmt.Printf("%s\tR(A)=%v B=%v C=%v\n", op, *a, bx, cx)

		case types.OP_JMP:
			// A sBx | pc+=sBx; if (A) close all upvalues >= R(A) + 1
			ax, bx := doJump(s, i, 0)
			fmt.Printf("%s\tA=%v sBx=%v\n", op, ax, bx)

		case types.OP_EQ, types.OP_LT, types.OP_LE:
			// A B C | if ((RK(B) == RK(C)) ~= A) then pc++
			// A B C | if ((RK(B) <  RK(C)) ~= A) then pc++
			// A B C | if ((RK(B) <= RK(C)) ~= A) then pc++
			ax := i.GetArgA()
			if _CMPOPS[op](*b, *c) != asBool(ax) {
				s.CI.PC++
			} else {
				// For the fall-through case, a JMP is always expected, in order to optimize
				// execution in the virtual machine.
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s\tA=%v RK(B)=%v RK(C)=%v\n", op, ax, *b, *c)

		case types.OP_TEST:
			// A C | if not (R(A) <=> C) then pc++
			cx, _ := i.GetArgC(false)
			if asBool(cx) == isFalse(*a) {
				s.CI.PC++
			} else {
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s\tR(A)=%v C=%v\n", op, *a, cx)

		case types.OP_TESTSET:
			// A B C | if (R(B) <=> C) then R(A) := R(B) else pc++
			cx, _ := i.GetArgC(false)
			if asBool(cx) == isFalse(*b) {
				s.CI.PC++
			} else {
				*a = *b
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s\tR(A)=%v R(B)=%v C=%v\n", op, *a, *b, cx)
			/*
				case types.OP_CALL, types.OP_TAILCALL: // TODO : For now, no tail call optimization
					ax := i.GetArgA()
					nParms, _ := i.GetArgB(false)
					nRets, _ := i.GetArgC(false)
					nRets--
					if nParms != 0 {
						// No parms: B=1, otherwise B-1 parms
						s.Stack.Top = s.CI.Base + ax + nParms
					}
					// Else, it is because last param to this call was a func call with unknown 
					// number of results, so this call actually set the Top to whatever it had to be.

					// TODO : See luaD_precall in ldo.c
					switch f := (*a).(type) {
					case types.GoFunc:
						// Go function call
						callGoFunc(s, f, s.CI.Base+ax+1, nRets)
					case *types.Closure:
						// Lune function call
						s.NewCallInfo(f, s.CI.Base+ax, nRets)
						goto newFrame
					}
					fmt.Printf("%s\n", op)

				case types.OP_RETURN:
					ax := s.CI.Base + i.GetArgA()
					bx, _ := i.GetArgB(false)
					if bx != 0 {
						s.Stack.Top = ax + bx - 1
					}
					if len(s.CI.Cl.P.Protos) > 0 {
						// TODO : Close upvalues
					}
					// TODO : See luaD_poscall in ldo.c, the hook magic is not implemented for now
					res := s.CI.FuncIndex
					wanted := s.CI.NumResults
					s.CI = s.CI.Prev
					// Set results in the right slots on the stack
					var j int
					for j = wanted; j != 0 && ax < s.Stack.Top; j-- {
						s.Stack.Stk[res] = s.Stack.Stk[ax]
						res, ax = res+1, ax+1
					}
					// Complete missing results with nils
					for ; j > 0; j-- {
						s.Stack.Stk[res] = nil
						res++
					}
					s.Stack.Top = res
					bx = wanted - types.LUNE_MULTRET
					if s.CI == nil {
						// TODO : Is this equivalent to Lua's check of CIST_REENTRY?
						fmt.Printf("%s\n", op)
						return
					} else {
						if bx != 0 {
							// TODO : Set Top back to CI.Top? I don't have a CI.Top!
						}
						if prevOp := s.CI.Cl.P.Code[s.CI.PC-1].GetOpCode(); prevOp != types.OP_CALL {
							panic(fmt.Sprintf("expected CALL to be previous instruction in RETURNed frame, got %s", prevOp))
						}
						fmt.Printf("%s : back to caller frame\n", op)
						goto newFrame
					}

				case types.OP_FORLOOP:
					ax := i.GetArgA()
					step := s.CI.Frame[ax+2].(float64)
					idx := s.CI.Frame[ax].(float64) + step
					limit := s.CI.Frame[ax+1].(float64)
					if 0 < step {
						if idx <= limit {
							s.CI.PC += i.GetArgsBx()
							s.CI.Frame[ax] = idx
							s.CI.Frame[ax+3] = idx
						}
					} else {
						if limit <= idx {
							s.CI.PC += i.GetArgsBx()
							s.CI.Frame[ax] = idx
							s.CI.Frame[ax+3] = idx
						}
					}
					fmt.Printf("%s : step:%v idx:%v limit:%v\n", op, step, idx, limit)

				case types.OP_FORPREP:
					ax := i.GetArgA()
					// TODO : Conversion, validate that it can be converted to a number (for now, assume number)
					init := (*a).(float64)
					limit := s.CI.Frame[ax+1].(float64)
					step := s.CI.Frame[ax+2].(float64)
					*a = init - step
					s.CI.PC += i.GetArgsBx()
					fmt.Printf("%s : init:%v limit:%v step:%v PC+=%v\n", op, init, limit, step, i.GetArgsBx())
			*/
		default:
			goto newFrame
			fmt.Printf("Ignore %s\n", op)
		}
	}
}
