package main

import (
	"fmt"
	"os"

	"github.com/mna/lune/serializer"
	"github.com/mna/lune/stdlib"
	"github.com/mna/lune/types"
	"github.com/mna/lune/vm"
)

func loadFile(fn string) (*types.Prototype, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return serializer.Load(f)
}

func main() {
	// Check args
	if len(os.Args) < 2 {
		fmt.Println("Expected an argument (file name)")
		os.Exit(1)
	}

	// Load file
	fn := os.Args[1]
	p, err := loadFile(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run file
	s := types.NewState(p)
	stdlib.OpenLibs(s.Globals)
	vm.Execute(s)
}
