library ieee;
   use ieee.std_logic_1164.all;

package p is
   --thdl:gen
   type t_rec is record
      slv : std_logic_vector(0 downto 0);
      suv : std_ulogic_vector(1 downto 0);
      si  : signed(2 downto 0);
      su  : unsigned(3 downto 0);
   end record;
end package;

package body p is
end package body;
