package stdlib

import (
	"github.com/PuerkitoBio/lune/types"
)

func OpenLibs(t types.Table) {
	libT := make(types.Table)
	libT[types.Value("write")] = nil
	t[types.Value("io")] = types.Value(libT)
}
