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

   --thdl:start
   -- Below code was automatically generated with the thdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-thdl.

   function to_rec(slv : std_logic_vector(9 downto 0)) return t_rec;
   function to_slv(rec : t_rec) return std_logic_vector;
   function to_str(rec : t_rec) return string;

   --thdl:end

end package;

package body p is

   --thdl:start
   -- Below code was automatically generated with the thdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-thdl.

   function to_rec(slv : std_logic_vector(9 downto 0)) return t_rec is
      variable rec : t_rec;
   begin
      rec.slv := slv(9 downto 9);
      rec.suv := slv(8 downto 7);
      rec.si := signed(slv(6 downto 4));
      rec.su := unsigned(slv(3 downto 0));
      return rec;
   end function;

   function to_slv(rec : t_rec) return std_logic_vector is
      variable slv : std_logic_vector(9 downto 0);
   begin
      slv(9 downto 9) := rec.slv;
      slv(8 downto 7) := rec.suv;
      slv(6 downto 4) := std_logic_vector(rec.si);
      slv(3 downto 0) := std_logic_vector(rec.su);
      return slv;
   end function;

   --thdl:end

end package body;
