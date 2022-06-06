library std;
   use std.textio.all;

library work;
   use work.pkg.all;

entity test is
end entity;

architecture tb of test is
begin
   main : process
      variable rec : t_rec := SOME_REC;
   begin
      write(output, to_str(rec) & LF);
      write(output, to_str(rec, true) & LF);

      std.env.finish;
   end process;
end architecture;
