package gen

import (
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/gen/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
	"strings"
	"sync"
)

var genArgs args.GenArgs

func Gen(args args.GenArgs) {
	genArgs = args

	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles := []string{}

	if genArgs.Filepath != "" {
		if strings.HasSuffix(genArgs.Filepath, ".vhd") || strings.HasSuffix(genArgs.Filepath, ".vhdl") {
			vhdlFiles = append(vhdlFiles, genArgs.Filepath)
		}
	} else {
		vhdlFiles = utils.GetVHDLFilePaths()
	}

	wg.Add(1)
	vhdl.Gen(args, vhdlFiles, &wg)
}
