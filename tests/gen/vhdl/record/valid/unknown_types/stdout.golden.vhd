library ieee;
   use ieee.std_logic_1164.all;

package p is
   --hdl:gen
   type t_foo is record
      f : t_field; --hdl: width=7
   end record;

   --hdl:gen
   type t_bar is record
      f : t_field; --hdl: width=4 to-type=lorem to-slv=ipsum to-str=dolor
   end record;

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_foo(slv : std_logic_vector(6 downto 0)) return t_foo;
   function to_slv(foo : t_foo) return std_logic_vector;
   function to_str(foo : t_foo; add_names : boolean := false) return string;

   function to_bar(slv : std_logic_vector(3 downto 0)) return t_bar;
   function to_slv(bar : t_bar) return std_logic_vector;
   function to_str(bar : t_bar; add_names : boolean := false) return string;

   --hdl:end

end package;

package body p is

   --hdl:start
   -- Below code was automatically generated with the hdl tool.
   -- Do not modify it by hand, unless you really know what you do.
   -- More info on https://github.com/m-kru/go-hdl.

   function to_foo(slv : std_logic_vector(6 downto 0)) return t_foo is
      variable foo : t_foo;
   begin
      foo.f := to_field(slv(6 downto 0));
      return foo;
   end function;

   function to_slv(foo : t_foo) return std_logic_vector is
      variable slv : std_logic_vector(6 downto 0);
   begin
      slv(6 downto 0) := to_slv(foo.f);
      return slv;
   end function;

   function to_str(foo : t_foo; add_names : boolean := false) return string is
   begin
      if add_names then
         return "(" & "f => " & to_str(foo.f) & ")";
      end if;
      return "(" & to_str(foo.f) & ")";
   end function;

   function to_bar(slv : std_logic_vector(3 downto 0)) return t_bar is
      variable bar : t_bar;
   begin
      bar.f := lorem(slv(3 downto 0));
      return bar;
   end function;

   function to_slv(bar : t_bar) return std_logic_vector is
      variable slv : std_logic_vector(3 downto 0);
   begin
      slv(3 downto 0) := ipsum(bar.f);
      return slv;
   end function;

   function to_str(bar : t_bar; add_names : boolean := false) return string is
   begin
      if add_names then
         return "(" & "f => " & dolor(bar.f) & ")";
      end if;
      return "(" & dolor(bar.f) & ")";
   end function;

   --hdl:end

end package body;
