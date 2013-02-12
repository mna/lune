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

func ioWrite(in []types.Value) []types.Value {
	ifs := make([]interface{}, len(in)+2)
	ifs = append(ifs, ">>")
	for _, v := range in {
		ifs = append(ifs, v)
	}
	ifs = append(ifs, "<<")
	fmt.Print(ifs...)
	return nil
}
