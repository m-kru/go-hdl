package args

var genHelpMsg string = `Gen command
===========

Usage
-----

  thdl gen [flags] [path/to/file]

Flags
  -to-stdout  Print to stdout instead of replacing file in place (useful for tests).

If path to file is not provided, thdl will scan all HDL files located in the tree
of working directory.


Description
-----------

The gen command scans HDL files and generates code based on their content.
The code generation is triggered by adding '--thdl:gen' tag line before symbol
for which the code should be generated. Note that there is no whitespace
between '--' and 'thdl:gen'. Inserting a whitespace in between is a good method
to temporarily disable code generation for a particular symbol. Another one is
adding an empty line between the '--thdl:gen' line and symbol line.

The code generation currently supports following kinds of symbols:

  VHDL:
    - enumeration type

        Example:
          --thdl:gen
          type t_status is (SUCCESS, ERROR);

        Thdl will generate following functions:
          - function to_status(slv : std_logic_vector(0 downto 0)) return t_status;
          - function to_slv(status : t_status) return std_logic_vector;
          - function to_str(status : t_status) return string;

        Parameters:
          - encoding  Encoding type. Valid encodings are: gray, one-hot, sequential.
                      The default encoding is sequential.

    - record type

        Example:
          --thdl:gen
          type t_data is record
             reverse : boolean;
             int     : integer;
             crc     : std_logic_vector(7 downto 0);
          end record;

        Thdl will generate following functions:
          - function to_data(slv : std_logic_vector(40 downto 0)) return t_data;
          - function to_slv(data : t_data) return std_logic_vector;
          - function to_str(data : t_data) return string;

        Flags:
          - no-to-str  Do not generate to_str function.

        By default thdl will handle all fields of standard types or of types that are also
        marked for generation and are located within the same scope (the scope is currently
        limited to the design unit). To handle foreign types one needs to provide additional
        information via field arguments. Field arguments are provided at the end of the line
        with particular field. They are prepended with '--thdl:' tag.
        Example:
          --thdl:gen
          type t_rec is record
             field : t_foreign; --thdl: width=N [to-type=name] [to-slv=name] [to-str=name]
          end record;
        Width is the only mandatory argument. N must be greater than 0. To-type is the name
        of the function used for the conversion from the std_logic_vector to the t_foreign.
        To-slv is the name of the function used for the conversion from the t_foreign to the
        std_logic_vector. To-str is the name of the function used for conversion from the
        t_foreign to the string. When names of these functions are not provided, then
        'to_foreign', 'to_slv' and 'to_str' are used respectively. Note that there are no
        whitespaces between argument name, its value and the '=' character.


Arguments passing
-----------------

To pass an argument to the 'thdl:gen' or 'thdl:' simply write parameter name
followed by the '=' character and actual argument value.

Examples:
  --thd:gen encoding=one-hot
  record_field : t_external_type; --thdl: width=8


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
