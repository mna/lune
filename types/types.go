package types

import (
	//"bytes"
	"fmt"
)

type ValType byte

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

var _ = fmt.Sprintf("")

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
	//var buf bytes.Buffer

	/*
			buf.WriteString(fmt.Sprintf("%+v\n", p.Meta))
			buf.WriteString(fmt.Sprintln("Instructions (", len(p.Code), ") :"))
			for _, c := range p.Code {
				buf.WriteString(fmt.Sprintln(c))
			}
			buf.WriteString(fmt.Sprintln("Constants (", len(p.Ks), ") :"))
			for _, v := range p.Ks {
				buf.WriteString(fmt.Sprintf("%+v\n", v))
			}
			buf.WriteString(fmt.Sprintln("Functions (", len(p.Protos), ") :"))
			for _, f := range p.Protos {
				buf.WriteString(fmt.Sprintln(f))
			}
			buf.WriteString(fmt.Sprintln("Upvalues (", len(p.Upvalues), ") :"))
			for _, u := range p.Upvalues {
				buf.WriteString(fmt.Sprintf("%+v\n", u))
			}
			buf.WriteString("\nDebug information:\n\n")
		buf.WriteString("Source: " + p.Source)
		buf.WriteString(fmt.Sprintln("Line info (", len(p.LineInfo), ") :"))
		buf.WriteString(fmt.Sprintln(p.LineInfo))
		buf.WriteString(fmt.Sprintln("Local variables (", len(p.LocVars), ") :"))
		for _, lv := range p.LocVars {
			buf.WriteString(fmt.Sprintf("%+v\n", lv))
		}*/

	//return buf.String()
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
