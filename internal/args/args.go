// Custom package for command line arguments parsing.
package args

import (
	"log"
	"strings"
)

type IgnoreList struct {
	ignore []string
}

func (il IgnoreList) FilterIgnored(filepaths []string) []string {
	ret := []string{}

	for _, fp := range filepaths {
		ignore := false
		for _, i := range il.ignore {
			if strings.Contains(fp, i) {
				ignore = true
				log.Printf("debug: ignoring %s\n", fp)
				break
			}
		}
		if !ignore {
			ret = append(ret, fp)
		}
	}

	return ret
}

type LibMap struct {
	libs map[string][]string
}

func (lm LibMap) Lib(filepath string) string {
	for name, patterns := range lm.libs {
		for _, p := range patterns {
			if strings.Contains(filepath, p) {
				return name
			}
		}
	}
	return ""
}

type VetArgs struct {
	IgnoreList
}

type DocArgs struct {
	IgnoreList
	LibMap
	Fusesoc    bool
	NoBold     bool
	SymbolPath string
}

type GenArgs struct {
	IgnoreList
	Fusesoc    bool
	NoBold     bool
	SymbolPath string
}

type Args struct {
	Cmd     string
	Debug   bool
	VetArgs VetArgs
	DocArgs DocArgs
	GenArgs GenArgs
}

func setFileCfgArgs(fc FileCfg, args *Args) {
	args.VetArgs.IgnoreList.ignore = fc.Vet.Ignore

	args.DocArgs.IgnoreList.ignore = fc.Doc.Ignore
	args.DocArgs.LibMap.libs = fc.Libs
	args.DocArgs.Fusesoc = fc.Doc.Fusesoc
	args.DocArgs.NoBold = fc.Doc.NoBold

	args.GenArgs.IgnoreList.ignore = fc.Gen.Ignore
}
