package args

var genHelpMsg string = `Gen command
===========

Usage
-----

  thdl gen [path/to/file]

If path to file is not provided, thdl will scan all HDL files located in the tree
of working directory.


Description
-----------

The gen command scans HDL files and generates code based on their content.
The code generation is triggered by adding '--thdl:gen' tag line before symbol
for which the code should be generated. Note that there is no whitespace
between '--' and 'thdl:gen'. Putting a whitespace in between is a good method
to temporarily disable code generation for particular symbol. Another one is
adding an empty line between the '--thdl:gen' line and symbol line.

The code generation currently supports following kinds of symbols:

  VHDL:
    - enumeration types

        Example:
          --thdl:gen
          type t_status is (SUCCESS, ERROR);

        Thdl will generate following functions:
          - function to_status(slv : std_logic_vector(0 downto 0)) return t_status;
          - function to_slv(status : t_status) return std_logic_vector;
          - function to_str(status : t_status) return string;


Naming symbols
--------------
It doesn't matter whether type symbol name has 't_' prefix.
For example, the names of the generated functions are the same for both
'type t_status is (A, B);' and 'type status is (A, B);'.


Constraints
-----------
File may contain multile design symbols, however package body must always
follow package declaration.
`
