package vm

type valtype uint

const (
	strTableCap = 100
)

const (
	vtNil valtype = iota
	vtBool
	vtNumber
	vtString
	vtFunction
	vtTable
	vtThread
	vtUserData // TODO : required?
)

/*
  Values are represented this way:
  nil:      value is nil
  bool:     value is bool
  number:   value is float64 or float32 based on GOARCH?
  string:   value is int, the index to the string table value
  function: value is int, the index to ?
  table:    ..
  thread:   ..
  userdata: ..
*/
type value interface{}

type table struct {
	m map[value]value
	a []value
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
	s []value
}

func newStack() *stack {
	return &stack{}
}
