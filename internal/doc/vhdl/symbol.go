package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/sym"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"os"
)

// symbol is a generic common symbol struct.
type symbol struct {
	parent sym.Symbol

	filepath string
	key      string
	name     string
	lineNum  uint32

	docStart uint32
	docEnd   uint32

	codeStart uint32
	codeEnd   uint32
}

func (s symbol) Filepath() string { return s.filepath }
func (s symbol) Key() string      { return s.key }
func (s symbol) Name() string     { return s.name }
func (s symbol) Files() []string  { panic("should never happen") }

func (s symbol) LineNum() uint32 { return s.lineNum }

func (s symbol) Doc() string {
	f, err := os.ReadFile(s.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", s.filepath, err)
	}

	return string(f[s.docStart:s.docEnd])
}

func (s symbol) Code() string {
	f, err := os.ReadFile(s.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", s.filepath, err)
	}

	return string(f[s.codeStart:s.codeEnd])
}

func (s symbol) DocCode() (string, string) {
	f, err := os.ReadFile(s.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", s.filepath, err)
	}

	doc := string(f[s.docStart:s.docEnd])
	code := string(f[s.codeStart:s.codeEnd])

	return doc, code
}

func (s symbol) OneLineSummary() string {
	code := utils.Dewhitespace(s.Code())
	if utils.IsSingleLine(code) {
		return fmt.Sprintf("%s", code)
	}
	return fmt.Sprintf("%s ...\n", utils.FirstLine(code))
}

func (s symbol) Path() string {
	return s.parent.Path() + "." + s.name
}
