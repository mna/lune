package stdlib

import (
	"github.com/PuerkitoBio/lune/types"
)

func OpenLibs(t types.Table) {
	var v, k1, k2, ti types.Value

	libT := make(types.Table)
	v = nil
	k1 = "write"
	libT[&k1] = &v
	k2 = "io"
	ti = libT
	t[&k2] = &ti
}
