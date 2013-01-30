package types

type ValType byte

const (
	strTableCap = 100
)

const (
	// Value types constants, must match with Lua's, see lua.h grep "basic types"
	TNIL ValType = iota
	TBOOL
	TLIGHTUSERDATA
	TNUMBER
	TSTRING
	TTABLE
	TFUNCTION
	TUSERDATA
	TTHREAD
)

/*
  Values are represented this way:
  nil:      value is nil
  bool:     value is bool
  number:   value is float64 or float32 based on GOARCH?
  string:   value is int, the index to the string table value
  function: value is *Prototype
  table:    ..
  thread:   ..
  userdata: ..
*/
type Value interface{}

type CallInfo struct {
	FuncIndex  int
	NumResults int
	CallStatus byte
	PC         int
	Base       int
}

// Naive implementation for now: always a map, no array optimization
type Table map[*Value]*Value

func (t *Table) Set(k *Value, v *Value) {
	t[k] = v
}

func (t *Table) Get(k *Value) *Value {
	return t[k]
}

/*
type table struct {
	m map[Value]Value
	a []Value
}

type strTable map[string]uint

func newStrTable() strTable {
	return make(strTable, strTableCap)
}

func (st strTable) intern(s string) bool {
	cnt, ok := st[s]
	if !ok {
		// This string doesn't exist yet
		st[s] = 1
	} else {
		// This string already exists, increment the refcount
		st[s] = cnt + 1
	}
	return ok
}

type stack struct {
	s []Value
}

func newStack() *stack {
	return &stack{}
}
*/
