package lune

type valtype uint

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
  nil:      v is ignored (is zero)
  bool:     v is either 0 (false) or 1 (true)
  number:   v is the actual value
  string:   v is the pointer to the string table value
  function: v is the pointer to ?
  table:    v is the pointer to ?
  thread:   v is the pointer to ?
  userdata: v is the pointer to ?
*/
type value struct {
	t valtype
	v float64
}
