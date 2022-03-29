  p_delay_proc : process (clk_i, rst_n_i)
  begin  -- process delay_proc
    if rst_n_i = '0' then               -- asynchronous reset (active low)
      genrst : for i in 1 to g_delay_cycles loop
        dly(i) <= (others => '0');
      end loop;
    elsif rising_edge(clk_i) then       -- rising clock edge
      dly(0) <= d_i;
      gendly : for i in 0 to g_delay_cycles-1 loop
        dly(i+1) <= dly(i);
      end loop;
    end if;
  end process p_delay_proc;
