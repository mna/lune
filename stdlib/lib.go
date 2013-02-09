package stdlib

import (
	"fmt"
	"github.com/PuerkitoBio/lune/types"
)

func OpenLibs(t types.Table) {
	var f types.GoFunc = ioWrite

	libT := make(types.Table)
	libT[types.Value("write")] = f
	t[types.Value("io")] = types.Value(libT)
}

func ioWrite(s *types.State) int {
	fmt.Println(">> In Go!")
	return 0
}
