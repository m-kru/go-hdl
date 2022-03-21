package args

var docHelpMsg string = `Doc command
===========

Usage
-----

thdl doc [flags] [<language>:][<library>.]<primarySymbol>[.<secondarySymbol>]

Flags:
  --no-bold  Don't print language keywords in bold.
  --fusesoc  Infer libraries from FuseSoc .core files.

Language is one of:
 - VHDL

Language is case insensitive.
Library, primarySymbol and secondarySymbol obey the case sensitivity of the language.

If language is omitted all languages are searched.
If multiple symbols are found, then ambiguity is reported.

If path to symbol consists of two words and a dot ("foo.bar"), then two scenarios
are assumed: library.primarySymbol and primarySymbol.secondarySymbol.
If multiple symbols are found, then ambiguity is reported.

If path to symbol is a single word, then two scenarios are assumed:
primarySymbol and secondarySymbol.
If multiple symbols are found, then ambiguity is reported.

Description
-----------

The doc command prints the documentation comment associated with the symbol
identified by its arguments followed by the source code of the symbol.

Following kinds of symbol are supported by the doc command:

  VHDL:
    - constant declaration within package declaration,
    - entity declaration,
    - function declaration within package declaration,
    - package declaration,
    - procedure declaration within package declaration,
    - type declaration within package declaration.
`
