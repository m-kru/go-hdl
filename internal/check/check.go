package check

import (
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/check/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"sync"
)

func Check(args args.CheckArgs) {
	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles, err := utils.GetFilePathsByExtension(".vhd", ".")
	if err != nil {
		log.Fatalf("%v", err)
	}
	wg.Add(1)
	vhdl.Check(vhdlFiles, &wg)
}
