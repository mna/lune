package vm

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/lune/types"
	"math"
)

func isFalse(v types.Value) bool {
	// Two values evaluate to False: nil and boolean false
	if v == nil {
		return true
	}
	if b, ok := v.(bool); ok && !b {
		return true
	}
	return false
}

func doJump(s *types.State, i types.Instruction, e int) (ax, bx int) {
	ax = i.GetArgA()
	if ax > 0 {
		// TODO : Close upvalues? See dojump in lvm.c.
	}
	bx = i.GetArgsBx()
	s.CI.PC += bx + e
	return
}

func callGoFunc(s *types.State, f types.GoFunc, base, nRets int) {
	var in []types.Value
	for i := base; i < s.Stack.Top; i++ {
		in = append(in, s.Stack.Get(i))
	}
	out := f(in)
	// Out values replace the stack values starting at the Go Func index (base - 1)
	// nRets values are expected, stop at this count, and fill with nils if necessary
	s.Stack.Top = base - 1
	for i := 0; i < nRets; i++ {
		if i < len(out) {
			s.Stack.Push(out[i])
		} else {
			s.Stack.Push(nil)
		}
	}
}

func Execute(s *types.State) {
	var a, b, c *types.Value

	// Start with entry point (position 0)
	s.NewCallInfo(s.Stack.Get(0).(*types.Closure), 0)

newFrame:
	for {
		i := s.CI.Cl.P.Code[s.CI.PC]
		op := i.GetOpCode()
		s.CI.PC++
		s.Stack.DumpStack()
		// TODO : Will fail big time here, when some opcodes are executed. OpMasks
		// just doesn't do what I think it does.
		a, b, c = i.GetArgs(s)

		switch op {
		case types.OP_MOVE:
			*a = *b
			fmt.Printf("%s : A=%v B=%v\n", i.GetOpCode(), *a, *b)

		case types.OP_LOADK:
			*a = *b
			fmt.Printf("%s : A=%v B=%v\n", op, *a, *b)

		case types.OP_LOADKx:
			// Special instruction: must always be followed by OP_EXTRAARG
			if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_EXTRAARG {
				panic(fmt.Sprintf("%s: expected OP_EXTRAARG as next instruction, found %s", op, i2.GetOpCode()))
			} else {
				s.CI.PC++
				ax := i2.GetArgAx()
				*a = s.CI.Cl.P.Ks[ax]
				fmt.Printf("%s : ax=%v a=%v\n", op, ax, *a)
			}

		case types.OP_LOADBOOL:
			bx, _ := i.GetArgB(false)
			bb := bx != 0
			*a = bb

			// Skip next instruction if C is true
			cx, _ := i.GetArgC(false)
			cb := cx != 0
			if cb {
				s.CI.PC++
			}
			fmt.Printf("%s : a=%v b=%v c=%v\n", op, *a, bb, cb)

		case types.OP_LOADNIL:
			ax := i.GetArgA()
			bx, _ := i.GetArgB(false)
			for j := 0; j <= bx; j++ {
				a = &s.CI.Frame[ax+j]
				*a = nil
			}
			fmt.Printf("%s : ax=%v bx=%v\n", op, ax, bx)

		case types.OP_GETUPVAL:
			*a = *b
			fmt.Printf("%s : A=%v B=%v\n", op, *a, *b)

		case types.OP_GETTABUP:
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s : k=%v v=%v ra=%v\n", op, *c, t.Get(*c), *a)

		case types.OP_GETTABLE:
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s : k=%v v=%v ra=%v\n", op, *c, t.Get(*c), *a)

		case types.OP_SETTABUP:
			t := (*a).(types.Table)
			t.Set(*b, *c)
			fmt.Printf("%s : k=%#v v=%#v\n", op, *b, *c)

		case types.OP_SETUPVAL:
			*b = *a
			fmt.Printf("%s : b=%v a=%v\n", op, *b, *a)

		case types.OP_SETTABLE:
			t := (*a).(types.Table)
			t.Set(*b, *c)
			fmt.Printf("%s : k=%v v=%v t=%v\n", op, *b, *c, *a)

		case types.OP_NEWTABLE:
			t := types.NewTable()
			// TODO : Encoded array and hash sizes (B and C) are ignored at the moment
			*a = t
			fmt.Printf("%s : t=%v b=%v c=%v\n", op, *a, *b, *c)

		case types.OP_SELF:
			ax := i.GetArgA()
			a = &s.CI.Frame[ax+1]
			*a = *b
			a = &s.CI.Frame[ax]
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s : a+1=%v a=%v c=%v\n", op, *b, t.Get(*c), *c)

			// TODO : For all operators, handle non-numeric data types (see Lua's conversion rules)
		case types.OP_ADD:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf + cf
			fmt.Printf("%s : b=%v + c=%v = a=%v\n", op, bf, cf, *a)

		case types.OP_SUB:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf - cf
			fmt.Printf("%s : b=%v - c=%v = a=%v\n", op, bf, cf, *a)

		case types.OP_MUL:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf * cf
			fmt.Printf("%s : b=%v * c=%v = a=%v\n", op, bf, cf, *a)

		case types.OP_DIV:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf / cf
			fmt.Printf("%s : b=%v ÷ c=%v = a=%v\n", op, bf, cf, *a)

		case types.OP_MOD:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = math.Mod(bf, cf)
			fmt.Printf("%s : b=%v % c=%v = a=%v\n", op, bf, cf, *a)

		case types.OP_POW:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = math.Pow(bf, cf)
			fmt.Printf("%s : b=%v ^ c=%v = a=%v\n", op, bf, cf, *a)

		case types.OP_UNM:
			bf := (*b).(float64)
			*a = -bf
			fmt.Printf("%s : -b=%v\n", op, *a)

		case types.OP_NOT:
			*a = isFalse(*b)
			fmt.Printf("%s : !b=%v\n", op, *a)

		case types.OP_LEN:
			switch v := (*b).(type) {
			case types.Table:
				*a = v.Len()
			case string:
				*a = len(v)
			default:
				// TODO : Metamethod	
			}
			fmt.Printf("%s : #b=%v\n", op, *a)

		case types.OP_CONCAT:
			// TODO : Manage type conversion? For now, assume strings...
			var buf bytes.Buffer
			bx, _ := i.GetArgB(false)
			cx, _ := i.GetArgC(false)
			for j := bx; j <= cx; j++ {
				buf.WriteString(s.CI.Frame[j].(string))
			}
			*a = buf.String()
			fmt.Printf("%s : b:%v .. c:%v = a:%v\n", op, *b, *c, *a)

		case types.OP_JMP:
			ax, bx := doJump(s, i, 0)
			fmt.Printf("%s : ax:%v PC+=%v\n", op, ax, bx)

		case types.OP_EQ:
			// Compares RK(B) and RK(C), which may be registers or constants. If the
			// boolean result is not A, then skip the next instruction. Conversely, if the
			// boolean result equals A, continue with the next instruction.
			ax := i.GetArgA()
			bf, cf := (*b).(float64), (*c).(float64)
			if (bf == cf) != (ax != 0) {
				s.CI.PC++
			} else {
				// For the fall-through case, a JMP is always expected, in order to optimize
				// execution in the virtual machine. In effect, EQ, LT and LE must always be
				// paired with a following JMP instruction.
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s : b:%v ==? c:%v | a:%v\n", op, bf, cf, ax)

		case types.OP_LT:
			// See OP_EQ for details.
			// TODO : See luaV_lessthan implementation, some subleties, type conversions?
			ax := i.GetArgA()
			bf, cf := (*b).(float64), (*c).(float64)
			if (bf < cf) != (ax != 0) {
				s.CI.PC++
			} else {
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s : b:%v <? c:%v | a:%v\n", op, bf, cf, ax)

		case types.OP_LE:
			// See OP_EQ for details.
			// TODO : See luaV_lessequal implementation, some subleties, type conversions?
			ax := i.GetArgA()
			bf, cf := (*b).(float64), (*c).(float64)
			if (bf <= cf) != (ax != 0) {
				s.CI.PC++
			} else {
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s : b:%v <=? c:%v | a:%v\n", op, bf, cf, ax)

		case types.OP_TEST:
			cx, _ := i.GetArgC(false)
			if (cx != 0) == isFalse(*a) {
				s.CI.PC++
			} else {
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s : c:%v !=? a:%v\n", op, cx, *a)

		case types.OP_TESTSET:
			cx, _ := i.GetArgC(false)
			if (cx != 0) == isFalse(*b) {
				s.CI.PC++
			} else {
				*a = *b
				if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_JMP {
					panic(fmt.Sprintf("%s: expected OP_JMP as next instruction, found %s", op, i2.GetOpCode()))
				} else {
					doJump(s, i2, 1)
				}
			}
			fmt.Printf("%s : c:%v !=? b:%v\n", op, cx, *b)

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
				s.NewCallInfo(f, s.CI.Base+ax)
				goto newFrame
			}
			fmt.Printf("%s\n", op)

		case types.OP_RETURN:
			if s.CI = s.CI.Prev; s.CI == nil {
				fmt.Printf("%s\n", op)
				return
			}
			fmt.Printf("%s : Back to previous call\n", op)
			goto newFrame

		default:
			fmt.Printf("Ignore %s\n", op)
		}

	}
	/*
	  CallInfo *ci = L->ci;
	  LClosure *cl;
	  TValue *k;
	  StkId base;
	 newframe:  // reentry point when frame changes (call/return) 
	  lua_assert(ci == L->ci);
	  cl = clLvalue(ci->func);
	  k = cl->p->k;
	  base = ci->u.l.base;
	  // main loop of interpreter 
	  for (;;) {
	    Instruction i = *(ci->u.l.savedpc++);
	    StkId ra;
	    if ((L->hookmask & (LUA_MASKLINE | LUA_MASKCOUNT)) &&
	        (--L->hookcount == 0 || L->hookmask & LUA_MASKLINE)) {
	      Protect(traceexec(L));
	    }
	    // WARNING: several calls may realloc the stack and invalidate `ra' 
	    ra = RA(i);
	    lua_assert(base == ci->u.l.base);
	    lua_assert(base <= L->top && L->top < L->stack + L->stacksize);
	    vmdispatch (GET_OPCODE(i)) {
	      vmcase(OP_TAILCALL,
	        int b = GETARG_B(i);
	        if (b != 0) L->top = ra+b;  // else previous instruction set top 
	        lua_assert(GETARG_C(i) - 1 == LUA_MULTRET);
	        if (luaD_precall(L, ra, LUA_MULTRET))  // C function? 
	          base = ci->u.l.base;
	        else {
	          // tail call: put called frame (n) in place of caller one (o) 
	          CallInfo *nci = L->ci;  // called frame 
	          CallInfo *oci = nci->previous;  // caller frame 
	          StkId nfunc = nci->func;  // called function 
	          StkId ofunc = oci->func;  // caller function 
	          // last stack slot filled by 'precall' 
	          StkId lim = nci->u.l.base + getproto(nfunc)->numparams;
	          int aux;
	          // close all upvalues from previous call 
	          if (cl->p->sizep > 0) luaF_close(L, oci->u.l.base);
	          // move new frame into old one 
	          for (aux = 0; nfunc + aux < lim; aux++)
	            setobjs2s(L, ofunc + aux, nfunc + aux);
	          oci->u.l.base = ofunc + (nci->u.l.base - nfunc);  // correct base 
	          oci->top = L->top = ofunc + (L->top - nfunc);  // correct top 
	          oci->u.l.savedpc = nci->u.l.savedpc;
	          oci->callstatus |= CIST_TAIL;  // function was tail called 
	          ci = L->ci = oci;  // remove new frame 
	          lua_assert(L->top == oci->u.l.base + getproto(ofunc)->maxstacksize);
	          goto newframe;  // restart luaV_execute over new Lua function 
	        }
	      )
	      vmcasenb(OP_RETURN,
	        int b = GETARG_B(i);
	        if (b != 0) L->top = ra+b-1;
	        if (cl->p->sizep > 0) luaF_close(L, base);
	        b = luaD_poscall(L, ra);
	        if (!(ci->callstatus & CIST_REENTRY))  // 'ci' still the called one 
	          return;  // external invocation: return 
	        else {  // invocation via reentry: continue execution 
	          ci = L->ci;
	          if (b) L->top = ci->top;
	          lua_assert(isLua(ci));
	          lua_assert(GET_OPCODE(*((ci)->u.l.savedpc - 1)) == OP_CALL);
	          goto newframe;  // restart luaV_execute over new Lua function 
	        }
	      )
	      vmcase(OP_FORLOOP,
	        lua_Number step = nvalue(ra+2);
	        lua_Number idx = luai_numadd(L, nvalue(ra), step); // increment index 
	        lua_Number limit = nvalue(ra+1);
	        if (luai_numlt(L, 0, step) ? luai_numle(L, idx, limit)
	                                   : luai_numle(L, limit, idx)) {
	          ci->u.l.savedpc += GETARG_sBx(i);  // jump back 
	          setnvalue(ra, idx);  // update internal index... 
	          setnvalue(ra+3, idx);  // ...and external index 
	        }
	      )
	      vmcase(OP_FORPREP,
	        const TValue *init = ra;
	        const TValue *plimit = ra+1;
	        const TValue *pstep = ra+2;
	        if (!tonumber(init, ra))
	          luaG_runerror(L, LUA_QL("for") " initial value must be a number");
	        else if (!tonumber(plimit, ra+1))
	          luaG_runerror(L, LUA_QL("for") " limit must be a number");
	        else if (!tonumber(pstep, ra+2))
	          luaG_runerror(L, LUA_QL("for") " step must be a number");
	        setnvalue(ra, luai_numsub(L, nvalue(ra), nvalue(pstep)));
	        ci->u.l.savedpc += GETARG_sBx(i);
	      )
	      vmcasenb(OP_TFORCALL,
	        StkId cb = ra + 3;  // call base 
	        setobjs2s(L, cb+2, ra+2);
	        setobjs2s(L, cb+1, ra+1);
	        setobjs2s(L, cb, ra);
	        L->top = cb + 3;  // func. + 2 args (state and index) 
	        Protect(luaD_call(L, cb, GETARG_C(i), 1));
	        L->top = ci->top;
	        i = *(ci->u.l.savedpc++);  // go to next instruction 
	        ra = RA(i);
	        lua_assert(GET_OPCODE(i) == OP_TFORLOOP);
	        goto l_tforloop;
	      )
	      vmcase(OP_TFORLOOP,
	        l_tforloop:
	        if (!ttisnil(ra + 1)) {  // continue loop? 
	          setobjs2s(L, ra, ra + 1);  // save control variable 
	           ci->u.l.savedpc += GETARG_sBx(i);  // jump back 
	        }
	      )
	      vmcase(OP_SETLIST,
	        int n = GETARG_B(i);
	        int c = GETARG_C(i);
	        int last;
	        Table *h;
	        if (n == 0) n = cast_int(L->top - ra) - 1;
	        if (c == 0) {
	          lua_assert(GET_OPCODE(*ci->u.l.savedpc) == OP_EXTRAARG);
	          c = GETARG_Ax(*ci->u.l.savedpc++);
	        }
	        luai_runtimecheck(L, ttistable(ra));
	        h = hvalue(ra);
	        last = ((c-1)*LFIELDS_PER_FLUSH) + n;
	        if (last > h->sizearray)  // needs more space? 
	          luaH_resizearray(L, h, last);  // pre-allocate it at once 
	        for (; n > 0; n--) {
	          TValue *val = ra+n;
	          luaH_setint(L, h, last--, val);
	          luaC_barrierback(L, obj2gco(h), val);
	        }
	        L->top = ci->top;  // correct top (in case of previous open call) 
	      )
	      vmcase(OP_CLOSURE,
	        Proto *p = cl->p->p[GETARG_Bx(i)];
	        Closure *ncl = getcached(p, cl->upvals, base);  // cached closure 
	        if (ncl == NULL)  // no match? 
	          pushclosure(L, p, cl->upvals, base, ra);  // create a new one 
	        else
	          setclLvalue(L, ra, ncl);  // push cashed closure 
	        checkGC(L, ra + 1);
	      )
	      vmcase(OP_VARARG,
	        int b = GETARG_B(i) - 1;
	        int j;
	        int n = cast_int(base - ci->func) - cl->p->numparams - 1;
	        if (b < 0) {  // B == 0? 
	          b = n;  // get all var. arguments 
	          Protect(luaD_checkstack(L, n));
	          ra = RA(i);  // previous call may change the stack 
	          L->top = ra + n;
	        }
	        for (j = 0; j < b; j++) {
	          if (j < n) {
	            setobjs2s(L, ra + j, base - n + j);
	          }
	          else {
	            setnilvalue(ra + j);
	          }
	        }
	      )
	      vmcase(OP_EXTRAARG,
	        lua_assert(0);
	      )
	    }
	  }*/
}
