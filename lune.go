package main

import (
	"fmt"
	"github.com/PuerkitoBio/lune/serializer"
	"github.com/PuerkitoBio/lune/stdlib"
	"github.com/PuerkitoBio/lune/vm"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected arguments (file name)")
		os.Exit(1)
	}
	fn := os.Args[1]
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	p, err := serializer.Load(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := vm.NewState(p)
	stdlib.OpenLibs(s.Globals)
	vm.Execute(s)
}
