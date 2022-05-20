package gen

import (
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/gen/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"strings"
	"sync"
)

var genArgs args.GenArgs

func Gen(args args.GenArgs) {
	genArgs = args

	var wg sync.WaitGroup
	defer wg.Wait()

	var err error
	vhdlFiles := []string{}

	if genArgs.Filepath != "" {
		if strings.HasSuffix(genArgs.Filepath, ".vhd") || strings.HasSuffix(genArgs.Filepath, ".vhdl") {
			vhdlFiles = append(vhdlFiles, genArgs.Filepath)
		}
	} else {
		vhdlFiles, err = utils.GetFilePathsByExtension(".vhd", ".")
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	wg.Add(1)
	vhdl.Gen(args, vhdlFiles, &wg)
}
