library ieee;
   use ieee.std_logic_1164.all;

package p1 is
   --hdl:gen
   type t_enum is (ZERO, ONE, TWO);
end package;

package body p1 is
end package body;

library ieee;
   use ieee.std_logic_1164.all;

package P2 is
   --hdl:gen encoding=one-hot
   type enum is (
      ZERO,
      one, Two,
      THREE
   );
end package;

package body p2 is
end P2;
