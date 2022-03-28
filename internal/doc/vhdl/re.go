package vhdl

import "regexp"

var emptyLine *regexp.Regexp = regexp.MustCompile(`^\s*$`)
var commentLine *regexp.Regexp = regexp.MustCompile(`^\s*--`)
var end *regexp.Regexp = regexp.MustCompile(`^\s*end\b`)
var endWithSemicolon *regexp.Regexp = regexp.MustCompile(`^\s*end\s*;`)
var endsWithSemicolon *regexp.Regexp = regexp.MustCompile(`;\s*($|--)`)

var constantDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*constant\s+(\w+)\s*(,\s*\w+)?\s*(,\s*\w+)?`)

var entityDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*entity\s+(\w*)\s+is`)
var architectureDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*architecture\s+\w+\s+of\s*\w+\s+is\b`)

var functionDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*(pure\b|impure\b)?\s*function\s+(\w+)`)
var endsWithReturn *regexp.Regexp = regexp.MustCompile(`\breturn\s+\w+\s*;`)

var packageDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*package\s+(\w*)\s+is`)
var endPackage *regexp.Regexp = regexp.MustCompile(`^\s*end\s+package\b`)
var packageBodyDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*package\s+body\s+\w+\s+is\b`)

var arrayTypeDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*type\s+(\w+)\s+is\s+array\b`)

var enumTypeDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*type\s+(\w+)\s+is\s*\(`)

var recordTypeDeclaration *regexp.Regexp = regexp.MustCompile(`^\s*type\s+(\w+)\s+is\s+record\b`)
var endRecord *regexp.Regexp = regexp.MustCompile(`^\s*end\s+record\b`)
