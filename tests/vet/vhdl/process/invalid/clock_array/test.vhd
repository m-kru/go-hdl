gen_rst_sync : for I in g_CLOCKS-1 downto 0 generate
  sync : process(clks_i, gc_reset_async_in)
  begin
    if gc_reset_async_in = '1' then
      rst_chains(i) <= (others => '0');
    elsif rising_edge(clks_i(i)) then
      rst_chains(i) <= '1' & rst_chains(i)(g_RST_LEN-1 downto 1);
    end if;
  end process;
  rst_n_o(i) <= rst_chains(i)(0);
end generate gen_rst_sync;
