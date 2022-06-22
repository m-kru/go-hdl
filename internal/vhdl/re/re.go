// Package re contains common regular expressions for VHDL.
package re

import re "regexp"

var EmptyLine *re.Regexp = re.MustCompile(`^\s*$`)
var CommentLine *re.Regexp = re.MustCompile(`^\s*--`)
var End *re.Regexp = re.MustCompile(`(?i)^\s*end\b`)
var EndWithSemicolon *re.Regexp = re.MustCompile(`(?i)^\s*end\s*;`)
var EndsWithSemicolon *re.Regexp = re.MustCompile(`(?i);\s*($|--)`)
var EndsWithRoundBracketAndSemicolon *re.Regexp = re.MustCompile(`(?i)\)\s*;\s*($|--)`)

var ConstantDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*constant\s+(\w+)\s*(,\s*\w+)?\s*(,\s*\w+)?`)

var EntityDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*entity\s+(\w*)\s+is`)
var ArchitectureDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*architecture\s+(\w+)\s+of\s*\w+\s+is\b`)

var FunctionDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*(pure\b|impure\b)?\s*function\s+(\w+)`)
var EndsWithReturn *re.Regexp = re.MustCompile(`(?i)\breturn\s+\w+\s*;`)

var PackageDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*package\s+(\w+)\s+is`)
var PackageInstantiation *re.Regexp = re.MustCompile(`(?i)^\s*package\s+(\w+)\s+is\s+new\b`)
var EndPackage *re.Regexp = re.MustCompile(`(?i)^\s*end\s+package\b`)
var PackageBodyDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*package\s+body\s+(\w+)\s+is\b`)
var EndPackageBody *re.Regexp = re.MustCompile(`(?i)^\s*end\s+package\s+body\b`)

var ProcedureDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*procedure\s+(\w+)`)

var ArrayTypeDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*type\s+(\w+)\s+is\s+array\b`)

var EnumTypeDeclaration *re.Regexp = re.MustCompile(`^(?i)\s*type\s+(\w+)\s+is\s*\(`)

var ProtectedTypeDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*type\s+(\w+)\s+is\s+protected\b`)
var EndProtected *re.Regexp = re.MustCompile(`(?i)^\s*end\s+protected\b`)

var RecordTypeDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*type\s+(\w+)\s+is\s+record\b`)
var EndRecord *re.Regexp = re.MustCompile(`(?i)^\s*end\s+record\b`)

var SubtypeDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*subtype\s+(\w+)\s+is\s+`)

var SomeTypeDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*type\s+(\w+)`)
var StartsWithArray *re.Regexp = re.MustCompile(`(?i)^\s*(is\s+)?array`)
var StartsWithProtected *re.Regexp = re.MustCompile(`(?i)^\s*(is\s+)?protected`)
var StartsWithRecord *re.Regexp = re.MustCompile(`(?i)^\s*(is\s+)?record`)
var StartsWithRoundBracket *re.Regexp = re.MustCompile(`(?i)^\s*(is\s+)?\(`)

var VariableDeclaration *re.Regexp = re.MustCompile(`(?i)^\s*(shared)?\s*variable\s+(\w+)\b`)

var SimpleRange *re.Regexp = re.MustCompile(`(?i)\s*(.+)\s+(downto|to)\s+(.+)\s*`)
