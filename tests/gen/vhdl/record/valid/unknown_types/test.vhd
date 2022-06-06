library ieee;
   use ieee.std_logic_1164.all;

package p is
   --thdl:gen
   type t_foo is record
      f : t_field; --thdl: width=7
   end record;

   --thdl:gen
   type t_bar is record
      f : t_field; --thdl: width=4 to-type=lorem to-slv=ipsum to-str=dolor
   end record;
end package;

package body p is
end package body;
