library ieee;
   use ieee.std_logic_1164.all;

package p is

   --thdl:gen
   type t_rec is record
      b    : bit;
      bool : boolean;
      sl   : std_logic;
      su   : std_ulogic;
   end record;

end package;

package body p is

end package body;
