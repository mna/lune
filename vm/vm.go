package vm

import (
	"github.com/PuerkitoBio/lune/serializer"
	"github.com/PuerkitoBio/lune/types"
)

func prepareFunction(s *State, fIdx int) (*serializer.Prototype, *types.CallInfo) {
	// Get the function's prototype at this stack index
	f := s.Get(fIdx)
	p := f.(*serializer.Prototype)

	// Make sure the stack has enough slots
	s.checkStack(p.Meta.MaxStackSize)

	// Complete the arguments
	n := s.stack.top - fIdx - 1
	for ; n < int(p.Meta.NumParams); n++ {
		s.push(nil)
	}

	ci := new(types.CallInfo)
	ci.FuncIndex = fIdx
	ci.NumResults = 0 // TODO : For now, ignore, someday will be passed
	ci.CallStatus = 0 // TODO : For now, ignore
	ci.PC = 0
	ci.Base = fIdx + 1 // TODO : For now, considre the base to be fIdx + 1, will have to manage varargs someday

	return p, ci

	// TODO : Var arg processing

	/*
	   case LUA_TLCL: {  // Lua function: prepare its call 
	     StkId base;
	     Proto *p = clLvalue(func)->p;
	     luaD_checkstack(L, p->maxstacksize);
	     func = restorestack(L, funcr);
	     n = cast_int(L->top - func) - 1;  // number of real arguments
	     for (; n < p->numparams; n++)
	       setnilvalue(L->top++);  // complete missing arguments
	     base = (!p->is_vararg) ? func + 1 : adjust_varargs(L, p, n);
	     ci = next_ci(L);  // now 'enter' new function
	     ci->nresults = nresults;
	     ci->func = func;
	     ci->u.l.base = base;
	     ci->top = base + p->maxstacksize;
	     lua_assert(ci->top <= L->stack_last);
	     ci->u.l.savedpc = p->code;  // starting point
	     ci->callstatus = CIST_LUA;
	     L->top = ci->top;
	     if (L->hookmask & LUA_MASKCALL)
	       callhook(L, ci);
	     return 0;
	   }
	*/
}

func Execute(s *State) {
	// Start with entry point (position 0)
	p, ci := prepareFunction(s, 0)

newFrame:
	for {
		i := p.Code[ci.PC]
		ci.PC++
		ra := ci.Base + i.GetArgA()

		switch i.GetOpCode() {
		case types.OP_SETTABUP:
      int a = GETARG_A(i);
      Protect(luaV_settable(L, cl->upvals[a]->v, RKB(i), RKC(i)));
		}

	}
	// TODO : CHeck bookmarks for how the CallInfo->u.l gets set (luaD_preCall)

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
	      vmcase(OP_MOVE,
	        setobjs2s(L, ra, RB(i));
	      )
	      vmcase(OP_LOADK,
	        TValue *rb = k + GETARG_Bx(i);
	        setobj2s(L, ra, rb);
	      )
	      vmcase(OP_LOADKX,
	        TValue *rb;
	        lua_assert(GET_OPCODE(*ci->u.l.savedpc) == OP_EXTRAARG);
	        rb = k + GETARG_Ax(*ci->u.l.savedpc++);
	        setobj2s(L, ra, rb);
	      )
	      vmcase(OP_LOADBOOL,
	        setbvalue(ra, GETARG_B(i));
	        if (GETARG_C(i)) ci->u.l.savedpc++;  // skip next instruction (if C) 
	      )
	      vmcase(OP_LOADNIL,
	        int b = GETARG_B(i);
	        do {
	          setnilvalue(ra++);
	        } while (b--);
	      )
	      vmcase(OP_GETUPVAL,
	        int b = GETARG_B(i);
	        setobj2s(L, ra, cl->upvals[b]->v);
	      )
	      vmcase(OP_GETTABUP,
	        int b = GETARG_B(i);
	        Protect(luaV_gettable(L, cl->upvals[b]->v, RKC(i), ra));
	      )
	      vmcase(OP_GETTABLE,
	        Protect(luaV_gettable(L, RB(i), RKC(i), ra));
	      )
	      vmcase(OP_SETTABUP,
	        int a = GETARG_A(i);
	        Protect(luaV_settable(L, cl->upvals[a]->v, RKB(i), RKC(i)));
	      )
	      vmcase(OP_SETUPVAL,
	        UpVal *uv = cl->upvals[GETARG_B(i)];
	        setobj(L, uv->v, ra);
	        luaC_barrier(L, uv, ra);
	      )
	      vmcase(OP_SETTABLE,
	        Protect(luaV_settable(L, ra, RKB(i), RKC(i)));
	      )
	      vmcase(OP_NEWTABLE,
	        int b = GETARG_B(i);
	        int c = GETARG_C(i);
	        Table *t = luaH_new(L);
	        sethvalue(L, ra, t);
	        if (b != 0 || c != 0)
	          luaH_resize(L, t, luaO_fb2int(b), luaO_fb2int(c));
	        checkGC(L, ra + 1);
	      )
	      vmcase(OP_SELF,
	        StkId rb = RB(i);
	        setobjs2s(L, ra+1, rb);
	        Protect(luaV_gettable(L, rb, RKC(i), ra));
	      )
	      vmcase(OP_ADD,
	        arith_op(luai_numadd, TM_ADD);
	      )
	      vmcase(OP_SUB,
	        arith_op(luai_numsub, TM_SUB);
	      )
	      vmcase(OP_MUL,
	        arith_op(luai_nummul, TM_MUL);
	      )
	      vmcase(OP_DIV,
	        arith_op(luai_numdiv, TM_DIV);
	      )
	      vmcase(OP_MOD,
	        arith_op(luai_nummod, TM_MOD);
	      )
	      vmcase(OP_POW,
	        arith_op(luai_numpow, TM_POW);
	      )
	      vmcase(OP_UNM,
	        TValue *rb = RB(i);
	        if (ttisnumber(rb)) {
	          lua_Number nb = nvalue(rb);
	          setnvalue(ra, luai_numunm(L, nb));
	        }
	        else {
	          Protect(luaV_arith(L, ra, rb, rb, TM_UNM));
	        }
	      )
	      vmcase(OP_NOT,
	        TValue *rb = RB(i);
	        int res = l_isfalse(rb);  // next assignment may change this value 
	        setbvalue(ra, res);
	      )
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
