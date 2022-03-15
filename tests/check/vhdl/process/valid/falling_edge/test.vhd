process (clk_i) is
begin
   if falling_edge(clk_i) then
      q <= d;
   end if;
end process;
