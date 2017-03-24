package vm

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mna/lune/types"
)

func asBool(i int) bool {
	return i != 0
}

func isNil(v types.Value) bool {
	return v == nil
}

func isFalse(v types.Value) bool {
	// Two values evaluate to False: nil and boolean false
	if isNil(v) {
		return true
	}
	if b, ok := v.(bool); ok && !b {
		return true
	}
	return false
}

func computeBinaryOp(op byte, b, c float64) float64 {
	switch op {
	case '+':
		return b + c
	case '-':
		return b - c
	case '*':
		return b * c
	case '/':
		return b / c
	case '%':
		return math.Mod(b, c)
	case '^':
		return math.Pow(b, c)
	}
	panic("unreachable")
}

func coerceAndComputeUnaryOp(op byte, b types.Value) types.Value {
	bf, bok := coerceToNumber(b)
	if bok {
		switch op {
		case '-':
			return -bf
		}
	} else {
		// TODO : Metamethods
		panic("metamethods not implemented")
	}
	return nil
}

func coerceAndComputeBinaryOp(op byte, b, c types.Value) types.Value {
	bf, bok := coerceToNumber(b)
	cf, cok := coerceToNumber(c)
	if bok && cok {
		// Both are numbers (or could be coerced to numbers)
		return computeBinaryOp(op, bf, cf)
	} else {
		// TODO : Metamethods
		panic("metamethods not implemented")
	}
	return nil
}

func coerceToNumber(v types.Value) (float64, bool) {
	switch bv := v.(type) {
	case float64:
		return bv, true
	case string:
		// Remove whitespace
		bv = strings.Trim(bv, " ")
		// First try to parse as an int
		if vi, err := strconv.ParseInt(bv, 0, 64); err == nil { // TODO : Int64 fits in float64?
			return float64(vi), true
		} else if vf, err := strconv.ParseFloat(bv, 64); err == nil { // TODO : Float64 for floats?
			return vf, true
		} else {
			return 0, false
		}
	default:
		return 0, false
	}
	panic("unreachable")
}

func coerceToString(v types.Value) (string, bool) {
	switch bv := v.(type) {
	case string:
		return bv, true
	case float64:
		// First try as an int
		vi := int64(bv)
		if float64(vi) == bv {
			return fmt.Sprintf("%d", vi), true
		} else {
			return fmt.Sprintf("%g", bv), true
		}
	default:
		return "", false
	}
	panic("unreachable")
}

func coerceAndConcatenate(src []types.Value) types.Value {
	var buf bytes.Buffer

	// Stop at i < len - 1 because the loop
	// uses i and i+1
	for i := 0; i < len(src); i++ {
		s, ok := coerceToString(src[i])
		if ok {
			buf.WriteString(s)
		} else {
			// TODO : Metamethods
			panic("metamethods not implemented")
		}
	}
	return buf.String()
}

func computeLength(v types.Value) float64 {
	switch bv := v.(type) {
	case types.Table:
		return float64(bv.Len())
	case string:
		return float64(len(bv))
	default:
		// TODO : Metamethod
		panic("metamethods not implemented")
	}
	panic("unreachable")
}

func areEqual(v1, v2 types.Value) bool {
	if t, ok := areSameType(v1, v2); !ok {
		return false
	} else if t == types.TNIL {
		return true
	}
	// TODO : Metamethods? No?
	return v1 == v2
}

func isLessEqual(l, r types.Value) bool {
	if t, ok := areSameType(l, r); ok {
		switch t {
		case types.TNUMBER:
			ln, rn := l.(float64), r.(float64)
			return ln <= rn
		case types.TSTRING:
			ls, rs := l.(string), r.(string)
			return ls <= rs
		}
	}
	// Not same type or not two numbers/strings
	// TODO : Metamethods, subtlety compared to LessThan, see lvm.c#luaV_lessequal
	return false
}

func isLessThan(l, r types.Value) bool {
	if t, ok := areSameType(l, r); ok {
		switch t {
		case types.TNUMBER:
			ln, rn := l.(float64), r.(float64)
			return ln < rn
		case types.TSTRING:
			ls, rs := l.(string), r.(string)
			return ls < rs
		}
	}
	// Not same type or not two numbers/strings
	// TODO : Metamethods
	return false
}

func areSameType(v1, v2 types.Value) (types.ValType, bool) {
	var valTypes [2]types.ValType
	var vals = [2]types.Value{v1, v2}

	for i := 0; i < 2; i++ {
		switch vals[i].(type) {
		case nil:
			valTypes[i] = types.TNIL
		case bool:
			valTypes[i] = types.TBOOL
		case float64:
			valTypes[i] = types.TNUMBER
		case string:
			valTypes[i] = types.TSTRING
		case *types.Closure:
			valTypes[i] = types.TFUNCTION
		case *types.Table:
			valTypes[i] = types.TTABLE
		}
	}
	return valTypes[0], valTypes[0] == valTypes[1]
}
