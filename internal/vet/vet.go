package vet

import (
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/utils"
	"github.com/m-kru/go-thdl/internal/vet/vhdl"
	"sync"
)

func Vet(args args.VetArgs) {
	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles := utils.GetVHDLFilePaths()
	wg.Add(1)
	vhdl.Vet(vhdlFiles, &wg)
}
