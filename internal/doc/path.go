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
}

func (sp symbolPath) String() string {
	return fmt.Sprintf(
		"Symbol path:\n  Language:  %s\n  Library:   %s\n  Primary:   %s\n  Secondary: %s\n",
		sp.language, sp.library, sp.primary, sp.secondary,
	)
}

// resolveSymbolPath returns a list of possible symbol paths based on the string.
func resolveSymbolPath(path string) []symbolPath {
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

	elems := strings.Split(innerPath, ".")

	// Temporary list of symbol paths before adding language.
	tmp := []symbolPath{}

	if len(elems) == 1 {
		tmp = append(tmp, symbolPath{primary: elems[0]})
		tmp = append(tmp, symbolPath{secondary: elems[0]})
	} else if len(elems) == 2 {
		tmp = append(tmp, symbolPath{library: elems[0], primary: elems[1]})
		tmp = append(tmp, symbolPath{primary: elems[0], secondary: elems[1]})
	} else if len(elems) == 3 {
		tmp = append(
			tmp,
			symbolPath{
				library:   elems[0],
				primary:   elems[1],
				secondary: elems[2],
			},
		)
	} else {
		log.Fatalf("invalid inner path format '%s', to many elements", innerPath)
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
