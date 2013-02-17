package vm

import (
	"github.com/PuerkitoBio/lune/types"
	"math"
	"strconv"
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
	}
	return nil
}

func coerceToNumber(v types.Value) (float64, bool) {
	switch bv := v.(type) {
	case float64:
		return bv, true
	case string:
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

func computeLength(v types.Value) float64 {
	switch bv := v.(type) {
	case types.Table:
		return float64(bv.Len())
	case string:
		return float64(len(bv))
	default:
		// TODO : Metamethod  
	}
	panic("unreachable")
}
