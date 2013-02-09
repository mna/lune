package vm

import (
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

func Execute(s *types.State) {
	var a, b, c *types.Value

	// Start with entry point (position 0)
	s.NewCallInfo(0, nil)

newFrame:
	for {
		i := s.CI.Cl.P.Code[s.CI.PC]
		s.CI.PC++
		s.Stack.DumpStack()
		a, b, c = i.GetArgs(s)

		switch i.GetOpCode() {
		case types.OP_MOVE:
			*a = *b
			fmt.Printf("%s : A=%v B=%v\n", i.GetOpCode(), *a, *b)

		case types.OP_LOADK:
			*a = *b
			fmt.Printf("%s : A=%v B=%v\n", i.GetOpCode(), *a, *b)

		case types.OP_LOADKx:
			// Special instruction: must always be followed by OP_EXTRAARG
			if i2 := s.CI.Cl.P.Code[s.CI.PC]; i2.GetOpCode() != types.OP_EXTRAARG {
				panic(fmt.Sprintf("%s: expected OP_EXTRAARG as next instruction, found %s", i.GetOpCode(), i2.GetOpCode()))
			} else {
				s.CI.PC++
				ax := i2.GetArgAx()
				*a = s.CI.Cl.P.Ks[ax]
				fmt.Printf("%s : ax=%v a=%v\n", i.GetOpCode(), ax, *a)
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
			fmt.Printf("%s : a=%v b=%v c=%v\n", i.GetOpCode(), *a, bb, cb)

		case types.OP_LOADNIL:
			ax := i.GetArgA()
			bx, _ := i.GetArgB(false)
			for j := 0; j <= bx; j++ {
				a = &s.CI.Frame[ax+j]
				*a = nil
			}
			fmt.Printf("%s : ax=%v bx=%v\n", i.GetOpCode(), ax, bx)

		case types.OP_GETUPVAL:
			*a = *b
			fmt.Printf("%s : A=%v B=%v\n", i.GetOpCode(), *a, *b)

		case types.OP_GETTABUP:
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s : k=%v v=%v ra=%v\n", i.GetOpCode(), *c, t.Get(*c), *a)

		case types.OP_GETTABLE:
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s : k=%v v=%v ra=%v\n", i.GetOpCode(), *c, t.Get(*c), *a)

		case types.OP_SETTABUP:
			t := (*a).(types.Table)
			t.Set(*b, *c)
			fmt.Printf("%s : k=%#v v=%#v\n", i.GetOpCode(), *b, *c)

		case types.OP_SETUPVAL:
			*b = *a
			fmt.Printf("%s : b=%v a=%v\n", i.GetOpCode(), *b, *a)

		case types.OP_SETTABLE:
			t := (*a).(types.Table)
			t.Set(*b, *c)
			fmt.Printf("%s : k=%v v=%v t=%v\n", i.GetOpCode(), *b, *c, *a)

		case types.OP_NEWTABLE:
			t := types.NewTable()
			// TODO : Encoded array and hash sizes (B and C) are ignored at the moment
			*a = t
			fmt.Printf("%s : t=%v b=%v c=%v\n", i.GetOpCode(), *a, *b, *c)

		case types.OP_SELF:
			ax := i.GetArgA()
			a = &s.CI.Frame[ax+1]
			*a = *b
			a = &s.CI.Frame[ax]
			t := (*b).(types.Table)
			*a = t.Get(*c)
			fmt.Printf("%s : a+1=%v a=%v c=%v\n", i.GetOpCode(), *b, t.Get(*c), *c)

			// TODO : For all operators, handle non-numeric data types (see Lua's conversion rules)
		case types.OP_ADD:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf + cf
			fmt.Printf("%s : b=%v + c=%v = a=%v\n", i.GetOpCode(), bf, cf, *a)

		case types.OP_SUB:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf - cf
			fmt.Printf("%s : b=%v - c=%v = a=%v\n", i.GetOpCode(), bf, cf, *a)

		case types.OP_MUL:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf * cf
			fmt.Printf("%s : b=%v * c=%v = a=%v\n", i.GetOpCode(), bf, cf, *a)

		case types.OP_DIV:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = bf / cf
			fmt.Printf("%s : b=%v ÷ c=%v = a=%v\n", i.GetOpCode(), bf, cf, *a)

		case types.OP_MOD:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = math.Mod(bf, cf)
			fmt.Printf("%s : b=%v % c=%v = a=%v\n", i.GetOpCode(), bf, cf, *a)

		case types.OP_POW:
			bf := (*b).(float64)
			cf := (*c).(float64)
			*a = math.Pow(bf, cf)
			fmt.Printf("%s : b=%v ^ c=%v = a=%v\n", i.GetOpCode(), bf, cf, *a)

		case types.OP_UNM:
			bf := (*b).(float64)
			*a = -bf
			fmt.Printf("%s : -b=%v\n", i.GetOpCode(), *a)

		case types.OP_NOT:
			*a = isFalse(*b)
			fmt.Printf("%s : !b=%v\n", i.GetOpCode(), *a)

		case types.OP_CALL:
			/*
				CALL A B C R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
				Performs a function call, with register R(A) holding the reference to the
				function object to be called. Parameters to the function are placed in the
				registers following R(A). If B is 1, the function has no parameters. If B is 2
				or more, there are (B-1) parameters.
				If B is 0, the function parameters range from R(A+1) to the top of the stack.
				This form is used when the last expression in the parameter list is a
				function call, so the number of actual parameters is indeterminate.
				Results returned by the function call is placed in a range of registers
				starting from R(A). If C is 1, no return results are saved. If C is 2 or more,
				(C-1) return values are saved. If C is 0, then multiple return results are
				saved, depending on the called function.
				CALL always updates the top of stack value. CALL, RETURN, VARARG
				and SETLIST can use multiple values (up to the top of the stack.)
			*/
			if f, ok := (*a).(types.GoFunc); ok {
				n := f(s)
				fmt.Printf("%s : %d\n", i.GetOpCode(), n)
			} else {
				fmt.Printf("%s : Ignored as not a GoFunc, not implemented yet.\n", i.GetOpCode())
			}

		case types.OP_RETURN:
			if s.CI = s.CI.Prev; s.CI == nil {
				fmt.Printf("%s\n", i.GetOpCode())
				return
			}
			fmt.Printf("%s : Back to previous call\n", i.GetOpCode())
			goto newFrame

		default:
			fmt.Printf("Ignore %s\n", i.GetOpCode())
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
	      vmcase(OP_LEN,
	        Protect(luaV_objlen(L, ra, RB(i)));
	      )
	      vmcase(OP_CONCAT,
	        int b = GETARG_B(i);
	        int c = GETARG_C(i);
	        StkId rb;
	        L->top = base + c + 1;  // mark the end of concat operands 
	        Protect(luaV_concat(L, c - b + 1));
	        ra = RA(i);  // 'luav_concat' may invoke TMs and move the stack 
	        rb = b + base;
	        setobjs2s(L, ra, rb);
	        checkGC(L, (ra >= rb ? ra + 1 : rb));
	        L->top = ci->top;  // restore top 
	      )
	      vmcase(OP_JMP,
	        dojump(ci, i, 0);
	      )
	      vmcase(OP_EQ,
	        TValue *rb = RKB(i);
	        TValue *rc = RKC(i);
	        Protect(
	          if (cast_int(equalobj(L, rb, rc)) != GETARG_A(i))
	            ci->u.l.savedpc++;
	          else
	            donextjump(ci);
	        )
	      )
	      vmcase(OP_LT,
	        Protect(
	          if (luaV_lessthan(L, RKB(i), RKC(i)) != GETARG_A(i))
	            ci->u.l.savedpc++;
	          else
	            donextjump(ci);
	        )
	      )
	      vmcase(OP_LE,
	        Protect(
	          if (luaV_lessequal(L, RKB(i), RKC(i)) != GETARG_A(i))
	            ci->u.l.savedpc++;
	          else
	            donextjump(ci);
	        )
	      )
	      vmcase(OP_TEST,
	        if (GETARG_C(i) ? l_isfalse(ra) : !l_isfalse(ra))
	            ci->u.l.savedpc++;
	          else
	          donextjump(ci);
	      )
	      vmcase(OP_TESTSET,
	        TValue *rb = RB(i);
	        if (GETARG_C(i) ? l_isfalse(rb) : !l_isfalse(rb))
	          ci->u.l.savedpc++;
	        else {
	          setobjs2s(L, ra, rb);
	          donextjump(ci);
	        }
	      )
	      vmcase(OP_CALL,
	        int b = GETARG_B(i);
	        int nresults = GETARG_C(i) - 1;
	        if (b != 0) L->top = ra+b;  // else previous instruction set top 
	        if (luaD_precall(L, ra, nresults)) {  // C function? 
	          if (nresults >= 0) L->top = ci->top;  // adjust results 
	          base = ci->u.l.base;
	        }
	        else {  // Lua function 
	          ci = L->ci;
	          ci->callstatus |= CIST_REENTRY;
	          goto newframe;  // restart luaV_execute over new Lua function 
	        }
	      )
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
