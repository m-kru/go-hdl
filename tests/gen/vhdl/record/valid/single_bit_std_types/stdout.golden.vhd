library ieee;
   use ieee.std_logic_1164.all;

package p is
   --thdl:gen
   type t_rec is record
      bi : bit;
      bo : boolean;
      sl : std_logic;
      su : std_ulogic;
   end record;

   --thdl:start
   -- Below code was automatically generated with the thdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-thdl.

   function to_rec(slv : std_logic_vector(3 downto 0)) return t_rec;
   function to_slv(rec : t_rec) return std_logic_vector;
   function to_str(rec : t_rec) return string;

   --thdl:end

end package;

package body p is

   --thdl:start
   -- Below code was automatically generated with the thdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-thdl.

   function to_rec(slv : std_logic_vector(3 downto 0)) return t_rec is
      variable rec : t_rec;
   begin
      if slv(3) = '1' then
         rec.bi := '1';
      elsif slv(3) = '0' then
         rec.bi := '0';
      else
         report "bit 3: cannot convert " & to_string(slv(3)) & " to bit type" severity failure;
      end if;
      if slv(2) = '1' then
         rec.bo := true;
      elsif slv(2) = '0' then
         rec.bo := false;
      else
         report "bit 2: cannot convert " & to_string(slv(2)) & " to boolean type" severity failure;
      end if;
      rec.sl := slv(1);
      rec.su := slv(0);
      return rec;
   end function;

   function to_slv(rec : t_rec) return std_logic_vector is
      variable slv : std_logic_vector(3 downto 0);
   begin
      if rec.bi = '1' then slv(3) := '1'; else slv(3) := '0'; end if;
      if rec.bo then slv(2) := '1'; else slv(2) := '0'; end if;
      slv(1) := rec.sl;
      slv(0) := rec.su;
      return slv;
   end function;

   --thdl:end

end package body;
