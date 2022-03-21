package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"log"
	"os"
)

type Entity struct {
	filepath string
	name     string

	hasDoc   bool
	docStart uint32
	docEnd   uint32

	codeStart uint32
	codeEnd   uint32
}

func (e Entity) Filepath() string { return e.filepath }

func (e Entity) Name() string { return e.name }

func (e Entity) Doc() string {
	f, err := os.ReadFile(e.filepath)
	if err != nil {
		log.Fatalf("reading '%s' entity code: error reading file %s: %v",
			e.name, e.filepath, err,
		)
	}

	return string(f[e.docStart:e.docEnd])
}

func (e Entity) Code() string {
	f, err := os.ReadFile(e.filepath)
	if err != nil {
		log.Fatalf("reading '%s' entity code: error reading file %s: %v",
			e.name, e.filepath, err,
		)
	}

	return string(f[e.codeStart:e.codeEnd])
}

func (e Entity) SymbolNames() []string {
	return []string{}
}

func (e Entity) GetSymbol(name string) (symbol.Symbol, bool) {
	return nil, false
}
