package vet

import (
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/utils"
	"github.com/m-kru/go-thdl/internal/vet/vhdl"
	"log"
	"sync"
)

func Vet(args args.VetArgs) {
	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles, err := utils.GetFilePathsByExtension(".vhd", ".")
	if err != nil {
		log.Fatalf("%v", err)
	}
	wg.Add(1)
	vhdl.Vet(vhdlFiles, &wg)
}
