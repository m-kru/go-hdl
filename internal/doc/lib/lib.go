package lib

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"os"
	"strings"
	"sync"
)

type LibrarySummary func(l *Library) string

type Library struct {
	lang string
	name string

	docFile    string
	files      []string
	filesMutex sync.Mutex

	symbols      map[string]symbol.Symbol
	symbolsMutex sync.Mutex

	libSummary LibrarySummary
}

func (l *Library) Filepath() string {
	return l.docFile
}

func (l *Library) Files() []string {
	return l.files
}

func (l *Library) LineNum() uint32 {
	panic("should never happen")
}

func MakeLibrary(lang string, name string, ls LibrarySummary) Library {
	if !utils.IsValidLang(lang) {
		panic("invalid language")
	}

	return Library{
		lang:       lang,
		name:       name,
		files:      []string{},
		symbols:    map[string]symbol.Symbol{},
		libSummary: ls,
	}
}

func (l *Library) AddFile(f string) {
	docFile := false

	switch l.lang {
	case "vhdl":
		if strings.HasSuffix(f, "doc.vhd") {
			docFile = true
		}
	default:
		panic("should never happen")
	}

	l.filesMutex.Lock()

	if docFile {
		if l.docFile != "" {
			log.Fatalf(
				"%s: library %s has at least 2 doc files:\n  %s\n  %s\n",
				l.lang, l.name, l.docFile, f,
			)
		}
		l.docFile = f
	}

	l.files = append(l.files, f)

	l.filesMutex.Unlock()
}

func (l *Library) Name() string { return l.name }

func (l *Library) SymbolNames() []string {
	names := []string{}

	for name, _ := range l.symbols {
		names = append(names, name)
	}

	return names
}

func (l *Library) Symbols() map[string]symbol.Symbol {
	return l.symbols
}

func (l *Library) GetSymbol(name string) []symbol.Symbol {
	if s, ok := l.symbols[name]; ok {
		return []symbol.Symbol{s}
	}
	return nil
}

func (l *Library) AddSymbol(s symbol.Symbol) {
	l.symbolsMutex.Lock()
	l.symbols[s.Name()] = s
	l.symbolsMutex.Unlock()
}

func (l *Library) Doc() string {
	if l.docFile == "" {
		return ""
	}

	f, err := os.ReadFile(l.docFile)
	if err != nil {
		log.Fatalf("error reading file %s: %v", l.docFile, err)
	}

	return string(f)
}

func (l *Library) Code() string {
	return l.libSummary(l)
}

func (l *Library) DocCode() (string, string) {
	return l.Doc(), l.Code()
}
