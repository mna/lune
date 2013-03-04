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
			[]types.Value{},
			types.Table{"a": 6.0},
			0,
		},
	}
)

func TestEnd2End(t *testing.T) {
	for _, tc := range end2endCases {
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
		vAc := s.Stack[i]
		if vEx != vAc {
			t.Errorf("%s: expected stack value %s at position %d, got %s", tc.name, vEx, i, vAc)
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
