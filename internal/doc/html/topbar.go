package html

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"strings"
)

func topbar(active string, nestingLevel int) string {
	homeActive := ""
	vhdlActive := ""

	switch active {
	case "home":
		homeActive = " active"
	case "vhdl":
		vhdlActive = " active"
	default:
		panic("should never happen")
	}

	root := "./"
	if nestingLevel == 1 {
		root = "../"
	} else if nestingLevel == 2 {
		root = "../../"
	} else if nestingLevel == 3 {
		root = "../../../"
	}

	b := strings.Builder{}

	b.WriteString(
		fmt.Sprintf("  <div class=\"topbar\">\n"+
			"    <div class=\"dropdown\">\n"+
			"      <button class=\"dropbtn%s\"><a href=\"%sindex.html\">Home</a></button>\n"+
			"    </div>\n", homeActive, root,
		),
	)

	vhdlLibs := vhdl.LibraryNames()
	if len(vhdlLibs) > 0 {
		b.WriteString(
			fmt.Sprintf("    <div class=\"dropdown\">\n"+
				"      <button class=\"dropbtn%s\"><a href=\"%svhdl/index.html\">VHDL</a></button>\n"+
				"      <div class=\"dropdown-content\">\n", vhdlActive, root,
			),
		)
		for _, l := range vhdlLibs {
			b.WriteString(
				fmt.Sprintf("        <a href=\"vhdl/%[1]s/index.html\">%[1]s</a>\n", l),
			)
		}
		b.WriteString(`      </div>
    </div>`)
	}

	b.WriteString(`  </div>`)

	return b.String()
}
