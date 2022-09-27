library ieee;
   use ieee.std_logic_1164.all;

package p1 is
   --hdl:gen
   type t_enum is (ZERO, ONE, TWO);

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_enum(slv : std_logic_vector(1 downto 0)) return t_enum;
   function to_slv(enum : t_enum) return std_logic_vector;
   function to_str(enum : t_enum) return string;

   --hdl:end

end package;

package body p1 is

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_enum(slv : std_logic_vector(1 downto 0)) return t_enum is
   begin
      case slv is
         when "00" => return ZERO;
         when "01" => return ONE;
         when "10" => return TWO;
         when others => report "invalid slv value " & to_string(slv) severity failure;
      end case;
   end function;

   function to_slv(enum : t_enum) return std_logic_vector is
   begin
      case enum is
         when ZERO => return "00";
         when ONE => return "01";
         when TWO => return "10";
      end case;
   end function;

   function to_str(enum : t_enum) return string is
   begin
      case enum is
         when ZERO => return "ZERO";
         when ONE => return "ONE";
         when TWO => return "TWO";
      end case;
   end function;

   --hdl:end

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

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_enum(slv : std_logic_vector(3 downto 0)) return enum;
   function to_slv(e : enum) return std_logic_vector;
   function to_str(e : enum) return string;

   --hdl:end

end package;

package body p2 is

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_enum(slv : std_logic_vector(3 downto 0)) return enum is
   begin
      case slv is
         when "0001" => return ZERO;
         when "0010" => return one;
         when "0100" => return Two;
         when "1000" => return THREE;
         when others => report "invalid slv value " & to_string(slv) severity failure;
      end case;
   end function;

   function to_slv(e : enum) return std_logic_vector is
   begin
      case e is
         when ZERO => return "0001";
         when one => return "0010";
         when Two => return "0100";
         when THREE => return "1000";
      end case;
   end function;

   function to_str(e : enum) return string is
   begin
      case e is
         when ZERO => return "ZERO";
         when one => return "one";
         when Two => return "Two";
         when THREE => return "THREE";
      end case;
   end function;

   --hdl:end

end P2;
