package vm

import (
	"fmt"
	"github.com/PuerkitoBio/lune/serializer"
	"github.com/PuerkitoBio/lune/types"
	"os"
	"testing"
)

type end2endTest struct {
	name    string
	opcodes []types.OpCode
	stack   []types.Value
	globals types.Table
	top     int
}

var (
	end2endCases = [...]end2endTest{
		end2endTest{
			"t1",
			[]types.OpCode{types.OP_SETTABUP, types.OP_RETURN},
			[]types.Value{nil},
			types.Table{"a": 6.0},
			0,
		},
		end2endTest{
			"t2",
			[]types.OpCode{types.OP_LOADK, types.OP_MUL, types.OP_SETTABUP, types.OP_RETURN},
			[]types.Value{nil, 10.5, 21.0},
			types.Table{"b": 21.0},
			0,
		},
		end2endTest{
			"t3",
			[]types.OpCode{types.OP_LOADK, types.OP_DIV, types.OP_TESTSET, types.OP_SUB, types.OP_RETURN},
			[]types.Value{nil, 7.0, 3.5, 3.5},
			types.Table{},
			0,
		},
		end2endTest{
			"t4",
			[]types.OpCode{
				types.OP_LOADBOOL,
				types.OP_LOADNIL,
				types.OP_SETTABUP,
				types.OP_SETTABUP,
				types.OP_GETTABUP,
				types.OP_NOT,
				types.OP_SETTABUP,
				types.OP_RETURN,
			},
			[]types.Value{nil, false, nil},
			types.Table{
				"a": true,
				"b": false,
			},
			0,
		},
	}
)

func TestEnd2End(t *testing.T) {
	for _, tc := range end2endCases {
		fmt.Printf("%s: running...\n", tc.name)
		testEnd2EndCase(t, tc)
	}
}

func testEnd2EndCase(t *testing.T, tc end2endTest) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("%s: PANIC %s", tc.name, err)
		}
	}()

	s, err := loadTestCase(tc)
	if err != nil {
		t.Errorf("%s: %s", tc.name, err)
	} else {
		Execute(s)
		assertTestCase(t, tc, s)
	}
}

func assertTestCase(t *testing.T, tc end2endTest, s *types.State) {
	assertOpcodes(t, tc, s)
	assertStack(t, tc, s)
	assertGlobals(t, tc, s)
}

func assertStack(t *testing.T, tc end2endTest, s *types.State) {
	// Oops, cannot check stack size like this, initialized at a size of 10...
	/*if lEx, lAc := len(tc.stack), len(s.Stack); lEx != lAc {
		t.Errorf("%s: expected %d stack size, got %d", tc.name, lEx, lAc)
	} else {*/
	// Same stack size, check values
	for i, vEx := range tc.stack {
		if i == 0 && vEx == nil {
			// Ignore if expected is nil at position 0 (will be the startup function)
			continue
		}
		vAc := s.Stack[i]
		if vEx != vAc {
			t.Errorf("%s: expected stack value %v at position %d, got %v", tc.name, vEx, i, vAc)
		}
	}
	//}

	if tc.top != s.Top {
		t.Errorf("%s: expected %d top-of-stack value, got %d", tc.name, tc.top, s.Top)
	}
}

func assertGlobals(t *testing.T, tc end2endTest, s *types.State) {
	if lEx, lAc := len(tc.globals), s.Globals.Len(); lEx != lAc {
		t.Errorf("%s: expected %d globals table size, got %d", tc.name, lEx, lAc)
	} else {
		for kEx, vEx := range tc.globals {
			vAc, ok := s.Globals[kEx]
			if !ok {
				t.Errorf("%s: expected key %v to exist in globals table", tc.name, kEx)
			} else if vEx != vAc {
				t.Errorf("%s: expected key %v in globals table to be %v, got %v", tc.name, kEx, vEx, vAc)
			}
		}
		for kAc, _ := range s.Globals {
			if _, ok := tc.globals[kAc]; !ok {
				t.Errorf("%s: found unexpected key %v in globals table", tc.name, kAc)
			}
		}
	}
}

func assertOpcodes(t *testing.T, tc end2endTest, s *types.State) {
	if lEx, lAc := len(tc.opcodes), len(s.OpCodeDebug); lEx != lAc {
		t.Errorf("%s: expected %d opcodes executed, got %d", tc.name, lEx, lAc)
	} else {
		// Same size, check values
		for i, opEx := range tc.opcodes {
			opAc := s.OpCodeDebug[i]
			if opEx != opAc {
				t.Errorf("%s: expected opcode %s at position %d, got %s", tc.name, opEx, i, opAc)
			}
		}
	}
}

func loadTestCase(tc end2endTest) (*types.State, error) {
	f, err := os.Open(fmt.Sprintf("./testdata/%s.out", tc.name))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p, err := serializer.Load(f)
	if err != nil {
		return nil, err
	}
	return types.NewState(p), nil
}
