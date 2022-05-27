library ieee;
   use ieee.std_logic_1164.all;

package p is
   --thdl:gen
   type t_enum is (ONE, TWO, THREE);
   --thdl:gen no-to-str
   type t_rec is record
      e : t_enum;
   end record;

   --thdl:start
   -- Below code was automatically generated with the thdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-thdl.

   function to_enum(slv : std_logic_vector(1 downto 0)) return t_enum;
   function to_slv(enum : t_enum) return std_logic_vector;
   function to_str(enum : t_enum) return string;

   function to_rec(slv : std_logic_vector(1 downto 0)) return t_rec;
   function to_slv(rec : t_rec) return std_logic_vector;

   --thdl:end

end package;

package body p is

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

   function to_rec(slv : std_logic_vector(1 downto 0)) return t_rec is
      variable rec : t_rec;
   begin
      rec.e := to_enum(slv(1 downto 0));
      return rec;
   end function;

   function to_slv(rec : t_rec) return std_logic_vector is
      variable slv : std_logic_vector(1 downto 0);
   begin
      slv(1 downto 0) := to_slv(rec.e);
      return slv;
   end function;

   --thdl:end

end package body;
