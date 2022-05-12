package args

import (
	"fmt"
	"strings"
)

type FileCfg struct {
	Ignore []string
	Libs   map[string][]string
	Vet    struct {
		Ignore []string
	}
	Doc struct {
		Ignore  []string
		Fusesoc bool
		NoBold  bool `yaml:"no-bold"`
		HTML    struct {
			Copyright string
			Path      string
			Title     string
		}
	}
	Gen struct {
		Ignore []string
	}
}

func (fc FileCfg) String() string {
	s := strings.Builder{}

	s.WriteString(".thdl.yml file configuration\n")
	s.WriteString("  Ignore:\n")
	for _, i := range fc.Ignore {
		s.WriteString(fmt.Sprintf("    - %s\n", i))
	}

	s.WriteString("  Libs:\n")
	for name, paths := range fc.Libs {
		s.WriteString(fmt.Sprintf("    %s\n", name))
		for _, p := range paths {
			s.WriteString(fmt.Sprintf("      - %s\n", p))
		}
	}

	s.WriteString("  Vet:\n")
	s.WriteString("    Ignore:\n")
	for _, i := range fc.Vet.Ignore {
		s.WriteString(fmt.Sprintf("      - %s\n", i))
	}

	s.WriteString("  Doc:\n")
	s.WriteString(fmt.Sprintf("    Fusesoc: %t\n", fc.Doc.Fusesoc))
	s.WriteString(fmt.Sprintf("    No-Bold: %t\n", fc.Doc.NoBold))
	s.WriteString("    HTML:\n")
	s.WriteString(fmt.Sprintf("      Copyright: %s\n", fc.Doc.HTML.Copyright))
	s.WriteString(fmt.Sprintf("      Path: %s\n", fc.Doc.HTML.Path))
	s.WriteString(fmt.Sprintf("      Title: %s\n", fc.Doc.HTML.Title))
	s.WriteString("    Ignore:\n")
	for _, i := range fc.Doc.Ignore {
		s.WriteString(fmt.Sprintf("      - %s\n", i))
	}

	s.WriteString("  Gen:\n")
	s.WriteString("    Ignore:\n")
	for _, i := range fc.Gen.Ignore {
		s.WriteString(fmt.Sprintf("      - %s\n", i))
	}

	return s.String()
}

func (fc *FileCfg) propagateGlobalIgnore() {
	for _, i := range fc.Ignore {
		fc.Vet.Ignore = append(fc.Vet.Ignore, i)
		fc.Doc.Ignore = append(fc.Doc.Ignore, i)
		fc.Gen.Ignore = append(fc.Gen.Ignore, i)
	}
}
