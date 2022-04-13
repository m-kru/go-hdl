package vhdl

import re "regexp"

var emptyLine *re.Regexp = re.MustCompile(`^\s*$`)
var commentLine *re.Regexp = re.MustCompile(`^\s*--`)
var end *re.Regexp = re.MustCompile(`^\s*end\b`)
var endWithSemicolon *re.Regexp = re.MustCompile(`^\s*end\s*;`)
var endsWithSemicolon *re.Regexp = re.MustCompile(`;\s*($|--)`)
var endsWithRoundBracketAndSemicolon *re.Regexp = re.MustCompile(`\)\s*;\s*($|--)`)

var constantDeclaration *re.Regexp = re.MustCompile(`^\s*constant\s+(\w+)\s*(,\s*\w+)?\s*(,\s*\w+)?`)

var entityDeclaration *re.Regexp = re.MustCompile(`^\s*entity\s+(\w*)\s+is`)
var architectureDeclaration *re.Regexp = re.MustCompile(`^\s*architecture\s+\w+\s+of\s*\w+\s+is\b`)

var functionDeclaration *re.Regexp = re.MustCompile(`^\s*(pure\b|impure\b)?\s*function\s+(\w+)`)
var endsWithReturn *re.Regexp = re.MustCompile(`\breturn\s+\w+\s*;`)

var packageDeclaration *re.Regexp = re.MustCompile(`^\s*package\s+(\w+)\s+is`)
var packageInstantiation *re.Regexp = re.MustCompile(`^\s*package\s+(\w+)\s+is\s+new\b`)
var endPackage *re.Regexp = re.MustCompile(`^\s*end\s+package\b`)
var packageBodyDeclaration *re.Regexp = re.MustCompile(`^\s*package\s+body\s+\w+\s+is\b`)

var procedureDeclaration *re.Regexp = re.MustCompile(`^\s*procedure\s+(\w+)`)

var arrayTypeDeclaration *re.Regexp = re.MustCompile(`^\s*type\s+(\w+)\s+is\s+array\b`)

var enumTypeDeclaration *re.Regexp = re.MustCompile(`^\s*type\s+(\w+)\s+is\s*\(`)

var protectedTypeDeclaration *re.Regexp = re.MustCompile(`^\s*type\s+(\w+)\s+is\s+protected\b`)
var endProtected *re.Regexp = re.MustCompile(`^\s*end\s+protected\b`)

var recordTypeDeclaration *re.Regexp = re.MustCompile(`^\s*type\s+(\w+)\s+is\s+record\b`)
var endRecord *re.Regexp = re.MustCompile(`^\s*end\s+record\b`)

var subtypeDeclaration *re.Regexp = re.MustCompile(`^\s*subtype\s+(\w+)\s+is\s+`)

var someTypeDeclaration *re.Regexp = re.MustCompile(`^\s*type\s+(\w+)`)
var startsWithArray *re.Regexp = re.MustCompile(`^\s*(is\s+)?array`)
var startsWithProtected *re.Regexp = re.MustCompile(`^\s*(is\s+)?protected`)
var startsWithRecord *re.Regexp = re.MustCompile(`^\s*(is\s+)?record`)
var startsWithRoundBracket *re.Regexp = re.MustCompile(`^\s*(is\s+)?\(`)
