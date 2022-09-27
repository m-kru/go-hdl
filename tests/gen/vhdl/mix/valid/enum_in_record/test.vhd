library ieee;
   use ieee.std_logic_1164.all;

package p is
   --hdl:gen
   type t_enum is (ONE, TWO, THREE);
   --hdl:gen no-to-str
   type t_rec is record
      e : t_enum;
   end record;
end package;

package body p is
end package body;
