package doc

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"strings"
)

type symbolPath struct {
	language  string
	library   string
	primary   string
	secondary string
	tertiary  string
}

func (sp symbolPath) String() string {
	var s string

	if sp.primary == "" {
		s = fmt.Sprintf("%s:%s", sp.language, sp.library)
	} else if sp.secondary == "" {
		s = fmt.Sprintf("%s:%s.%s", sp.language, sp.library, sp.primary)
	} else if sp.tertiary == "" {
		s = fmt.Sprintf("%s:%s.%s.%s", sp.language, sp.library, sp.primary, sp.secondary)
	} else {
		s = fmt.Sprintf(
			"%s:%s.%s.%s.%s",
			sp.language, sp.library, sp.primary, sp.secondary, sp.tertiary,
		)
	}
	return s
}

func (sp symbolPath) DebugString() string {
	return fmt.Sprintf(
		"{Language: %s, Library: %s, Primary: %s, Secondary: %s, Tertiary: %s}",
		sp.language, sp.library, sp.primary, sp.secondary, sp.tertiary,
	)
}

func (sp symbolPath) isLibrary() bool {
	if sp.primary == "" && sp.secondary == "" && sp.tertiary == "" {
		if sp.library != "" {
			return true
		}
		panic("should never happen")
	}
	return false
}

// resolveSymbolPath returns a list of possible symbol paths based on the string.
func resolveSymbolPath(path string) []symbolPath {
	if utils.IsTooGeneralPath(path) {
		log.Fatalf("provided path is too general")
	}

	langs := []string{}

	innerPath := path

	if strings.Contains(path, ":") {
		elems := strings.Split(path, ":")
		lang := strings.ToLower(elems[0])
		if !utils.IsValidLang(lang) {
			log.Fatalf("invalid language '%s'", lang)
		}
		langs = append(langs, lang)
		innerPath = elems[1]
	} else {
		for _, l := range utils.ValidLangs() {
			langs = append(langs, l)
		}
	}

	// Temporary list of symbol paths before adding language.
	var tmp []symbolPath

	// Unequivocal path.
	if innerPath[len(innerPath)-1] == '.' {
		innerPath = innerPath[0 : len(innerPath)-1]
		elems := strings.Split(innerPath, ".")

		switch len(elems) {
		case 1:
			tmp = []symbolPath{
				symbolPath{library: elems[0]},
			}
		case 2:
			tmp = []symbolPath{
				symbolPath{library: elems[0], primary: elems[1]},
			}
		case 3:
			tmp = []symbolPath{
				symbolPath{library: elems[0], primary: elems[1], secondary: elems[2]},
			}
		case 4:
			tmp = []symbolPath{
				symbolPath{
					library: elems[0], primary: elems[1], secondary: elems[2], tertiary: elems[3],
				},
			}
		default:
			log.Fatalf("invalid inner path format '%s', to many elements", innerPath)
		}
	} else {
		elems := strings.Split(innerPath, ".")

		switch len(elems) {
		case 1:
			tmp = []symbolPath{
				symbolPath{library: elems[0]},
				symbolPath{library: "*", primary: elems[0]},
				symbolPath{library: "*", primary: "*", secondary: elems[0]},
			}
		case 2:
			tmp = []symbolPath{
				symbolPath{library: elems[0], primary: elems[1]},
				symbolPath{library: "*", primary: elems[0], secondary: elems[1]},
				symbolPath{library: "*", primary: "*", secondary: elems[0], tertiary: elems[1]},
			}
		case 3:
			tmp = []symbolPath{
				symbolPath{library: elems[0], primary: elems[1], secondary: elems[2]},
				symbolPath{library: "*", primary: elems[0], secondary: elems[1], tertiary: elems[2]},
			}
		case 4:
			tmp = []symbolPath{
				symbolPath{library: elems[0], primary: elems[1], secondary: elems[2], tertiary: elems[3]},
			}
		default:
			log.Fatalf("invalid inner path format '%s', to many elements", innerPath)
		}
	}

	sps := []symbolPath{}

	for _, sp := range tmp {
		for _, l := range langs {
			sp.language = l
			if l == "vhdl" {
				sp.library = strings.ToLower(sp.library)
				sp.primary = strings.ToLower(sp.primary)
				sp.secondary = strings.ToLower(sp.secondary)
			}
			sps = append(sps, sp)
		}
	}

	return sps
}
