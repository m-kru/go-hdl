package lib

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
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

	symbols      map[string][]sym.Symbol
	symbolsMutex sync.Mutex

	libSummary LibrarySummary
}

func (l *Library) Filepath() string          { return l.docFile }
func (l *Library) Files() []string           { return l.files }
func (l *Library) Key() string               { return l.name }
func (l *Library) Name() string              { return l.name }
func (l *Library) Code() string              { return l.libSummary(l) }
func (l *Library) DocCode() (string, string) { return l.Doc(), l.Code() }
func (l *Library) OneLineSummary() string    { panic("not yet implemented") }
func (l *Library) LineNum() uint32           { panic("should never happen") }

func MakeLibrary(lang string, name string, ls LibrarySummary) *Library {
	if !utils.IsValidLang(lang) {
		panic("invalid language")
	}

	return &Library{
		lang:       lang,
		name:       name,
		files:      []string{},
		symbols:    map[string][]sym.Symbol{},
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

func (l *Library) InnerKeys() []string {
	keys := []string{}

	for key, _ := range l.symbols {
		keys = append(keys, key)
	}

	return keys
}

func (l *Library) Symbols() map[string][]sym.Symbol {
	return l.symbols
}

func (l *Library) GetSymbol(key string) []sym.Symbol {
	if s, ok := l.symbols[key]; ok {
		return s
	}
	return nil
}

func (l *Library) AddSymbol(s sym.Symbol) {
	l.symbolsMutex.Lock()

	if _, ok := l.symbols[s.Key()]; ok {
		l.symbols[s.Key()] = append(l.symbols[s.Key()], s)
	} else {
		l.symbols[s.Key()] = []sym.Symbol{s}
	}

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

func (l *Library) Path() string {
	return l.lang + ":" + l.name
}
