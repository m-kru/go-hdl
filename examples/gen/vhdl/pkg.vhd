library ieee;
   use ieee.std_logic_1164.all;
   use ieee.numeric_std.all;

package pkg is
   --thdl:gen
   type t_enum is (ONE, TWO, THREE);

   --thdl:gen
   type t_rec is record
      enum : t_enum;
      bi   : bit;
      bool : boolean;
      sl   : std_logic;
      slv  : std_logic_vector(7 downto 0);
      int  : integer;
   end record;

   constant SOME_REC : t_rec := (ONE, '0', false, 'X', "11000011", 123);

   --thdl:start
   -- Below code was automatically generated with the thdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-thdl.

   function to_enum(slv : std_logic_vector(1 downto 0)) return t_enum;
   function to_slv(enum : t_enum) return std_logic_vector;
   function to_str(enum : t_enum) return string;

   function to_rec(slv : std_logic_vector(44 downto 0)) return t_rec;
   function to_slv(rec : t_rec) return std_logic_vector;
   function to_str(rec : t_rec; add_names : boolean := false) return string;

   --thdl:end

end package;

package body pkg is

   --thdl:start
   -- Below code was automatically generated with the thdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-thdl.

   function to_enum(slv : std_logic_vector(1 downto 0)) return t_enum is
   begin
      case slv is
         when "00" => return ONE;
         when "01" => return TWO;
         when "10" => return THREE;
         when others => report "invalid slv value " & to_string(slv) severity failure;
      end case;
   end function;

   function to_slv(enum : t_enum) return std_logic_vector is
   begin
      case enum is
         when ONE => return "00";
         when TWO => return "01";
         when THREE => return "10";
      end case;
   end function;

   function to_str(enum : t_enum) return string is
   begin
      case enum is
         when ONE => return "ONE";
         when TWO => return "TWO";
         when THREE => return "THREE";
      end case;
   end function;

   function to_rec(slv : std_logic_vector(44 downto 0)) return t_rec is
      variable rec : t_rec;
   begin
      rec.enum := to_enum(slv(44 downto 43));
      if slv(42) = '1' then
         rec.bi := '1';
      elsif slv(42) = '0' then
         rec.bi := '0';
      else
         report "bit 42: cannot convert " & to_string(slv(42)) & " to bit type" severity failure;
      end if;
      if slv(41) = '1' then
         rec.bool := true;
      elsif slv(41) = '0' then
         rec.bool := false;
      else
         report "bit 41: cannot convert " & to_string(slv(41)) & " to boolean type" severity failure;
      end if;
      rec.sl := slv(40);
      rec.slv := slv(39 downto 32);
      rec.int := to_integer(signed(slv(31 downto 0)));
      return rec;
   end function;

   function to_slv(rec : t_rec) return std_logic_vector is
      variable slv : std_logic_vector(44 downto 0);
   begin
      slv(44 downto 43) := to_slv(rec.enum);
      if rec.bi = '1' then slv(42) := '1'; else slv(42) := '0'; end if;
      if rec.bool then slv(41) := '1'; else slv(41) := '0'; end if;
      slv(40) := rec.sl;
      slv(39 downto 32) := rec.slv;
      slv(31 downto 0) := std_logic_vector(to_signed(rec.int, 32));
      return slv;
   end function;
   function to_str(rec : t_rec; add_names : boolean := false) return string is
   begin
      if add_names then
         return "(" &"enum => " & to_str(rec.enum) & ", " &"bi => " & to_string(rec.bi) & ", " &"bool => " & to_string(rec.bool) & ", " &"sl => " & to_string(rec.sl) & ", " &"slv => " & to_string(rec.slv) & ", " &"int => " & to_string(rec.int) & ")";
      end if;
      return "(" & to_str(rec.enum) & ", " & to_string(rec.bi) & ", " & to_string(rec.bool) & ", " & to_string(rec.sl) & ", " & to_string(rec.slv) & ", " & to_string(rec.int) & ")";
   end function;

   --thdl:end

end package body;
