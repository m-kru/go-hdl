library ieee;
   use ieee.std_logic_1164.all;

package p is
   --hdl:gen
   type t_rec is record
      bi : bit;
      bo : boolean;
      sl : std_logic;
      su : std_ulogic;
   end record;

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_rec(slv : std_logic_vector(3 downto 0)) return t_rec;
   function to_slv(rec : t_rec) return std_logic_vector;
   function to_str(rec : t_rec; add_names : boolean := false) return string;

   --hdl:end

end package;

package body p is

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

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

   function to_str(rec : t_rec; add_names : boolean := false) return string is
   begin
      if add_names then
         return "(" & "bi => " & to_string(rec.bi) & ", " & "bo => " & to_string(rec.bo) & ", " & "sl => " & to_string(rec.sl) & ", " & "su => " & to_string(rec.su) & ")";
      end if;
      return "(" & to_string(rec.bi) & ", " & to_string(rec.bo) & ", " & to_string(rec.sl) & ", " & to_string(rec.su) & ")";
   end function;

   --hdl:end

end package body;
