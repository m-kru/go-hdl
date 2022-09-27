library ieee;
   use ieee.std_logic_1164.all;
   use ieee.numeric_std.all;

package p is
   --hdl:gen
   type t_rec is record
      i : integer;
      n : natural;
      p : positive;
   end record;

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_rec(slv : std_logic_vector(95 downto 0)) return t_rec;
   function to_slv(rec : t_rec) return std_logic_vector;
   function to_str(rec : t_rec; add_names : boolean := false) return string;

   --hdl:end

end package;

package body p is

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_rec(slv : std_logic_vector(95 downto 0)) return t_rec is
      variable rec : t_rec;
   begin
      rec.i := to_integer(signed(slv(95 downto 64)));
      rec.n := to_integer(unsigned(slv(63 downto 32)));
      rec.p := to_integer(unsigned(slv(31 downto 0)));
      return rec;
   end function;

   function to_slv(rec : t_rec) return std_logic_vector is
      variable slv : std_logic_vector(95 downto 0);
   begin
      slv(95 downto 64) := std_logic_vector(to_signed(rec.i, 32));
      slv(63 downto 32) := std_logic_vector(to_unsigned(rec.n, 32));
      slv(31 downto 0) := std_logic_vector(to_unsigned(rec.p, 32));
      return slv;
   end function;

   function to_str(rec : t_rec; add_names : boolean := false) return string is
   begin
      if add_names then
         return "(" & "i => " & to_string(rec.i) & ", " & "n => " & to_string(rec.n) & ", " & "p => " & to_string(rec.p) & ")";
      end if;
      return "(" & to_string(rec.i) & ", " & to_string(rec.n) & ", " & to_string(rec.p) & ")";
   end function;

   --hdl:end

end package body;
