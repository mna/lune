package lune

type strTable map[string]uint32

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

type gState struct {
	strt strTable
}
