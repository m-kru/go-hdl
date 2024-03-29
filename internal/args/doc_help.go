package args

var docHelpMsg string = `Doc command
===========

Usage
-----

  hdl doc [flags] symbolPath

Flags:
  -debug      Print debug messages.
  -no-bold    Don't print language keywords in bold.
  -no-config  Don't read .hdl.yml config file.
  -fusesoc    Infer libraries from FuseSoc .core files.
  -html       Generate documentation in HTML format.

Parameters:
  -html-copyright  Copyright placed in the left bottom corner.
  -html-path       Output path for the generated HTML documentation.
  -html-title      HTML title.


Description
-----------

The doc command prints the documentation comment associated with the symbol,
followed by the source code of the symbol.

The following kinds of symbols are supported by the doc command:

  VHDL:
    - constant declaration within package declaration,
    - entity declaration,
    - function declaration within package declaration,
    - package declaration,
    - package instantiation within package declaration,
    - procedure declaration within package declaration,
    - type declaration within package declaration.


Symbol path
-----------

The symbol path has following structure:

  language:library.primarySymbol.secondarySymbol.tertiarySymbol

Language is one of:
 - VHDL
If language is omitted, all languages are searched.

Language is case insensitive. Library, primarySymbol, secondarySymbol and
tertiarySymbol obey the case sensitivity of the language. A tertiary symbol
can't be '*' wildcard.

Symbol path can consist of one to five words.
If symbol path is a single word three scenarios are assumed:
  - library,
  - primarySymbol,
  - secondarySymbol.
If symbol path consists of two words and a dot ("foo.bar") three scenarios
are assumed:
  - library.primarySymbol,
  - primarySymbol.secondarySymbol,
  - secondarySymbol.tertiarySymbol.
If symbol path consists of three words and two dots ("foo.bar.baz") two
scenarios are assumed:
  - library.primarySymbol.secondarySymbol,
  - primarySymbol.secondarySymbol.tertiarySymbol.

If multiple symbols are found ambiguity is reported.

To resolve symbol path ambiguity extend the path by adding the preceding symbol
name. If symbols with the same name exist at different levels,  resolve
ambiguity by adding '.' at the end. For example, let's assume there is "foo"
library, "foo" primarySymbol, "foo" secondarySymbol, and "foo" tertiarySymbol.
To see the documentation of "foo" library type:
  hdl doc foo.
To see the documentation of "foo" primarySymbol type:
  hdl doc foo.foo.
To see the documentation of "foo" secondarySymbol type:
  hdl doc foo.foo.foo.
To see the documentation of "foo" tertiarySymbol type:
  hdl doc foo.foo.foo.foo


Library documentation
---------------------

To document a library provide 'doc.<langugeExtension>' file within the library.
For example, to document a VHDL library provide 'doc.vhd' file. Each library
can have only one doc file. If more than one doc file is found per library,
then the error is reported.


HTML Style
----------
Generated HTML style can be customized via style.css file located in the
generated css directory. One can edit style.css or replace it with
completely different one. For example, if your company uses particular
colors and you generate documentation using CI/CD tool, then prepare custom
style.css file and overwrite the generated one with cp command.
`
