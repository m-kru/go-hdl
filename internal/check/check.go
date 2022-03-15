package check

import (
	"github.com/m-kru/go-thdl/internal/check/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"sync"
)

func Check(cmdLineArgs map[string]string) {
	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles, err := utils.GetFilePathsByExtension(".vhd", ".")
	if err != nil {
		log.Fatalf("%v", err)
	}
	wg.Add(1)
	vhdl.Check(vhdlFiles, &wg)
}
