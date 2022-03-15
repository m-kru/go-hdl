process (clk_i) is
begin
   if rising_edge(clk_i) then
      q <= d;
   end if;
end process;
