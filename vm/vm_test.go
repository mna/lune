package vm

import (
	"fmt"
	"github.com/PuerkitoBio/lune/serializer"
	"github.com/PuerkitoBio/lune/types"
	"os"
	"reflect"
	"testing"
)

// Definition of an end to end test case
type end2endTest struct {
	name    string
	context string
	opcodes []types.OpCode
	stack   []types.Value
	globals types.Table
	top     int
}

var (
	end2endCases = [...]end2endTest{
		end2endTest{
			"t1",
			"",
			[]types.OpCode{types.OP_SETTABUP, types.OP_RETURN},
			[]types.Value{nil},
			types.Table{"a": 6.0},
			0,
		},
		end2endTest{
			"t2",
			"",
			[]types.OpCode{types.OP_LOADK, types.OP_MUL, types.OP_SETTABUP, types.OP_RETURN},
			[]types.Value{nil, 10.5, 21.0},
			types.Table{"b": 21.0},
			0,
		},
		end2endTest{
			"t3",
			"",
			[]types.OpCode{types.OP_LOADK, types.OP_DIV, types.OP_TESTSET, types.OP_SUB, types.OP_RETURN},
			[]types.Value{nil, 7.0, 3.5, 3.5},
			types.Table{},
			0,
		},
		end2endTest{
			"t4",
			"",
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
		end2endTest{
			"t5",
			"",
			[]types.OpCode{
				types.OP_NEWTABLE,
				types.OP_SETTABUP,
				types.OP_GETTABUP,
				types.OP_SETTABLE,
				types.OP_GETTABUP,
				types.OP_GETTABLE,
				types.OP_MUL,
				types.OP_SETTABUP,
				types.OP_RETURN,
			},
			[]types.Value{nil, 12.0},
			types.Table{
				"a": types.Table{"test": 6.0},
				"b": 12.0,
			},
			0,
		},
		end2endTest{
			"t6",
			"",
			[]types.OpCode{
				types.OP_LOADK,
				types.OP_LEN,
				types.OP_RETURN,
			},
			[]types.Value{nil, "I come from down in the valley", 30.0},
			types.Table{},
			0,
		},
	}
)

// Run all end to end test cases
func TestEnd2End(t *testing.T) {
	for _, tc := range end2endCases {
		fmt.Printf("%s: running...\n", tc.name)
		testEnd2EndCase(t, tc)
	}
}

// Test a single end to end test case
func testEnd2EndCase(t *testing.T, tc end2endTest) {
	/*defer func() {
		if err := recover(); err != nil {
			t.Errorf("%s: PANIC %s", tc.name, err)
		}
	}()*/

	s, err := loadTestCase(tc)
	if err != nil {
		t.Errorf("%s: %s", tc.name, err)
	} else {
		Execute(s)
		assertTestCase(t, tc, s)
	}
}

// Assert the expected results for a test case
func assertTestCase(t *testing.T, tc end2endTest, s *types.State) {
	assertOpcodes(t, tc, s)
	assertStack(t, tc, s)
	assertGlobals(t, tc, s)
}

// Assert the stack of the state
func assertStack(t *testing.T, tc end2endTest, s *types.State) {
	// Oops, cannot check stack size like this, initialized at a size of 10...
	/*if lEx, lAc := len(tc.stack), len(s.Stack); lEx != lAc {
		t.Errorf("%s: expected %d stack size, got %d", tc.name, lEx, lAc)
	} else {*/
	tc.context = "stack"
	// Same stack size, check values
	for i, vEx := range tc.stack {
		if i == 0 && vEx == nil {
			// Ignore if expected is nil at position 0 (will be the startup function)
			continue
		}
		vAc := s.Stack[i]
		assertValues(t, tc, vEx, vAc)
	}
	//}

	if tc.top != s.Top {
		t.Errorf("%s: expected %d top-of-stack value, got %d", tc.name, tc.top, s.Top)
	}
}

// Assert the globals table
func assertGlobals(t *testing.T, tc end2endTest, s *types.State) {
	tc.context = "globals"
	assertTables(t, tc, tc.globals, s.Globals)
}

func assertValues(t *testing.T, tc end2endTest, vEx types.Value, vAc types.Value) {
	typeEx, typeAc := reflect.TypeOf(vEx), reflect.TypeOf(vAc)
	// From reflect package's doc for String(): To test for equality, compare the Types directly.
	if typeEx != typeAc {
		t.Errorf("%s: expected %s value to be of type %s, got type %s", tc.name, tc.context, typeEx, typeAc)
	} else {
		// Same type, compare value
		if typeEx != nil && typeEx.Kind() == reflect.Map {
			// Maps are uncomparable, must use assertTables
			assertTables(t, tc, vEx.(types.Table), vAc.(types.Table))
		} else if vEx != vAc {
			t.Errorf("%s: expected %s value to be %v, got %v", tc.name, tc.context, vEx, vAc)
		}
	}
}

func assertTables(t *testing.T, tc end2endTest, tEx types.Table, tAc types.Table) {
	if lEx, lAc := tEx.Len(), tAc.Len(); lEx != lAc {
		t.Errorf("%s: expected %s table size to be %d, got %d", tc.name, tc.context, lEx, lAc)
	} else {
		ori := tc.context
		for kEx, vEx := range tEx {
			vAc, ok := tAc[kEx]
			if !ok {
				t.Errorf("%s: expected %s key %v to exist in table", tc.name, tc.context, kEx)
			} else {
				tc.context = fmt.Sprintf("%s.%v", ori, kEx)
				assertValues(t, tc, vEx, vAc)
			}
		}

		// Now look for unexpected keys in actual table
		for kAc, _ := range tAc {
			if _, ok := tEx[kAc]; !ok {
				t.Errorf("%s: unexpected %s key %v in table", tc.name, tc.context, kAc)
			}
		}
	}
}

// Assert the executed opcodes
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

// Load a test case
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
