package vhdl

import (
	"log"
	"os"
)

// Symbol is a generic common symbol struct.
type Symbol struct {
	filepath string
	name     string

	hasDoc   bool
	docStart uint32
	docEnd   uint32

	codeStart uint32
	codeEnd   uint32
}

func (s Symbol) Filepath() string { return s.filepath }

func (s Symbol) Name() string { return s.name }

func (s Symbol) Doc() string {
	f, err := os.ReadFile(s.filepath)
	if err != nil {
		log.Fatalf("reading '%s' entity code: error reading file %s: %v",
			s.name, s.filepath, err,
		)
	}

	return string(f[s.docStart:s.docEnd])
}

func (s Symbol) Code() string {
	f, err := os.ReadFile(s.filepath)
	if err != nil {
		log.Fatalf("reading '%s' entity code: error reading file %s: %v",
			s.name, s.filepath, err,
		)
	}

	return string(f[s.codeStart:s.codeEnd])
}
