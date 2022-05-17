package vhdl

import re "regexp"

var thdlGenLine *re.Regexp = re.MustCompile(`\s*--thdl:gen\b`)

var enumTypeDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*type\s+(\w+)\s+is\s*\(`)
