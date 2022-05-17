package gen

import (
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/gen/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"sync"
)

var genArgs args.GenArgs

func Gen(args args.GenArgs) {
	genArgs = args

	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles, err := utils.GetFilePathsByExtension(".vhd", ".")
	if err != nil {
		log.Fatalf("%v", err)
	}
	wg.Add(1)
	vhdl.Gen(vhdlFiles, &wg)
}
