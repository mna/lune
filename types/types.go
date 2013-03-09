package types

import (
	"fmt"
)

type ValType byte

const (
	// Value types constants, must match with Lua's, see lua.h grep "basic types"
	TNIL ValType = iota
	TBOOL
	_ // TLIGHTUSERDATA, not implemented yet
	TNUMBER
	TSTRING
	TTABLE
	TFUNCTION
	_ // TUSERDATA, not implemented yet
	_ // TTHREAD, not implemented yet
)

// Go function type
type GoFunc func([]Value) []Value

/*
  Values are represented this way:
  nil:      value is nil
  bool:     value is bool
  number:   value is float64 (TODO: or float32 based on GOARCH?)
  string:   value is string
  function: value is *Closure (lune) or GoFunc (Go)
  table:    value is Table
  thread:   ..
  userdata: ..
*/
type Value interface{}

type Closure struct {
	P      *Prototype
	UpVals []Value
}

func NewClosure(p *Prototype) *Closure {
	return &Closure{p, make([]Value, len(p.Upvalues))}
}

func (cl *Closure) String() string {
	return cl.P.String()
}

// Naive implementation for now: always a map, no array optimization
type Table map[Value]Value

func NewTable() Table {
	return make(Table)
}

func (t Table) Set(k Value, v Value) {
	t[k] = v
}

func (t Table) Get(k Value) Value {
	return t[k]
}

func (t Table) Len() int {
	// TODO : This is not how the # (length operator) works in Lua, see
	// http://www.lua.org/manual/5.2/manual.html#3.4.6
	return len(t)
}

type Prototype struct {
	Meta     *FuncMeta
	Code     []Instruction
	Ks       []Value
	Protos   []*Prototype
	Upvalues []*Upvalue

	// Debug info, unavailable in release build
	Source   string
	LineInfo []int32
	LocVars  []*LocVar
}

func (p *Prototype) String() string {
	return fmt.Sprintf("%s.%d", p.Source, p.Meta.LineDefined)
}

type FuncMeta struct {
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVarArg        byte
	MaxStackSize    byte
}

type Upvalue struct {
	Name    string
	Instack byte
	Idx     byte
}

type LocVar struct {
	Name    string
	Startpc uint32
	Endpc   uint32
}
