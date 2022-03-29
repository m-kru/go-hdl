package vhdl

import (
	"testing"
)

func TestPositiveResetStuckToOne(t *testing.T) {
	var tests = []struct {
		line string
		msg  string
		ok   bool
	}{
		// Invalid mappings
		{line: "rst=>'1'", msg: "positive reset stuck to '1'", ok: false},
		{line: "rstp => '1',", msg: "positive reset stuck to '1'", ok: false},
		{line: "rst_p=> '1',", msg: "positive reset stuck to '1'", ok: false},
		{line: "rst_i_p  => '1',", msg: "positive reset stuck to '1'", ok: false},
		{line: "resetp_i =>  '1',", msg: "positive reset stuck to '1'", ok: false},
		// Valid mappings
		{line: "rstp => '0',", msg: "", ok: true},
		{line: "reset => '0'", msg: "", ok: true},
	}

	for i, test := range tests {
		msg, ok := checkResetPortMapping([]byte(test.line))
		if msg != test.msg || ok != test.ok {
			t.Errorf("[%d]: got (%v, %v); want (%v, %v)", i, msg, ok, test.msg, test.ok)
		}
	}
}

func TestPositiveResetMappedToNegativeReset(t *testing.T) {
	var tests = []struct {
		line string
		msg  string
		ok   bool
	}{
		// Invalid mappings
		{line: "rst=>rstn", msg: "positive reset mapped to negative reset", ok: false},
		{line: "resetp => rst_n_i", msg: "positive reset mapped to negative reset", ok: false},
		// Valid mappings
		{line: "rstp => not rstn,", msg: "", ok: true},
		{line: "rst => rstp,", msg: "", ok: true},
		{line: "reset => not(rstn)", msg: "", ok: true},
		{line: "rst => rst_in", msg: "", ok: true},
		{line: "rst_i => rst_i", msg: "", ok: true},
		{line: "when c_RST => r_rst_n <= '0';", msg: "", ok: true},
	}

	for i, test := range tests {
		msg, ok := checkResetPortMapping([]byte(test.line))
		if msg != test.msg || ok != test.ok {
			t.Errorf("[%d]: got (%v, %v); want (%v, %v)", i, msg, ok, test.msg, test.ok)
		}
	}
}

func TestNegativeResetStuckToZero(t *testing.T) {
	var tests = []struct {
		line string
		msg  string
		ok   bool
	}{
		// Invalid mappings
		{line: "rstn=>'0'", msg: "negative reset stuck to '0'", ok: false},
		{line: "rst_n => '0',", msg: "negative reset stuck to '0'", ok: false},
		{line: "rst_i_n  => '0',", msg: "negative reset stuck to '0'", ok: false},
		{line: "resetn_i =>  '0',", msg: "negative reset stuck to '0'", ok: false},
		// Valid mappings
		{line: "rstn => '1',", msg: "", ok: true},
		{line: "reset_n => '1'", msg: "", ok: true},
	}

	for i, test := range tests {
		msg, ok := checkResetPortMapping([]byte(test.line))
		if msg != test.msg || ok != test.ok {
			t.Errorf("[%d]: got (%v, %v); want (%v, %v)", i, msg, ok, test.msg, test.ok)
		}
	}
}

func TestNegativeResetMappedToPositiveReset(t *testing.T) {
	var tests = []struct {
		line string
		msg  string
		ok   bool
	}{
		// Invalid mappings
		{line: "rstn=>rst", msg: "negative reset mapped to positive reset", ok: false},
		{line: "resetn => rst_i", msg: "negative reset mapped to positive reset", ok: false},
		// Valid mappings
		{line: "rstn => rstn,", msg: "", ok: true},
		{line: "rst_n_i => not rstp,", msg: "", ok: true},
		{line: "resetn => not(rstp)", msg: "", ok: true},
		{line: "rst_n => not rst_in", msg: "", ok: true},
		{line: "rst_i_n => not(rst_i)", msg: "", ok: true},
	}

	for i, test := range tests {
		msg, ok := checkResetPortMapping([]byte(test.line))
		if msg != test.msg || ok != test.ok {
			t.Errorf("[%d]: got (%v, %v); want (%v, %v)", i, msg, ok, test.msg, test.ok)
		}
	}
}

func TestPositiveResetIfCondition(t *testing.T) {
	var tests = []struct {
		line string
		msg  string
		ok   bool
	}{
		// Invalid mappings
		{line: "if rst='0' then", msg: "invalid positive reset condition", ok: false},
		{line: "if(rst='0')then", msg: "invalid positive reset condition", ok: false},
		{line: "   if rst_i = '0'  then ", msg: "invalid positive reset condition", ok: false},
		{line: "	if  resetp_i ='0'  then ", msg: "invalid positive reset condition", ok: false},
		{line: "if reset_p = '0'then", msg: "invalid positive reset condition", ok: false},
		{line: "if rst_i_p = '0'then", msg: "invalid positive reset condition", ok: false},
		{line: "if not rst then", msg: "invalid positive reset condition", ok: false},
		{line: "if not(rst_i) then", msg: "invalid positive reset condition", ok: false},
		{line: "if not ( reset ) then", msg: "invalid positive reset condition", ok: false},
		// Valid mappings
		{line: "if reset_p = '1' then", msg: "", ok: true},
		{line: "if (reset_p = '1') then", msg: "", ok: true},
		{line: "if rst_i = '1' then", msg: "", ok: true},
		{line: "   if reset ='1' then ", msg: "", ok: true},
		{line: "if reset then ", msg: "", ok: true},
		{line: "   if  rst_p_i  then ", msg: "", ok: true},
	}

	for i, test := range tests {
		msg, ok := checkPositiveResetIfCondition([]byte(test.line))
		if msg != test.msg || ok != test.ok {
			t.Errorf("[%d]: got (%v, %v); want (%v, %v)", i, msg, ok, test.msg, test.ok)
		}
	}
}

func TestNegativeResetIfCondition(t *testing.T) {
	var tests = []struct {
		line string
		msg  string
		ok   bool
	}{
		// Invalid mappings
		{line: "if rst_n ='1' then", msg: "invalid negative reset condition", ok: false},
		{line: "if (rst_n ='1') then", msg: "invalid negative reset condition", ok: false},
		{line: "if(rst_n ='1')then", msg: "invalid negative reset condition", ok: false},
		{line: "   if rst_i_n = '1'  then ", msg: "invalid negative reset condition", ok: false},
		{line: "	if  resetn_i ='1'  then ", msg: "invalid negative reset condition", ok: false},
		{line: "if reset_n = '1'then", msg: "invalid negative reset condition", ok: false},
		{line: "if rst_i_n = '1'then", msg: "invalid negative reset condition", ok: false},
		{line: "if rst_n then", msg: "invalid negative reset condition", ok: false},
		{line: "if rstn_i then", msg: "invalid negative reset condition", ok: false},
		{line: "if  ( reset_n ) then", msg: "invalid negative reset condition", ok: false},
		// Valid mappings
		{line: "if reset_n = '0' then", msg: "", ok: true},
		{line: "if rst_n_i = '0' then", msg: "", ok: true},
		{line: "   if reset_n ='0' then ", msg: "", ok: true},
		{line: "if not(reset_n) then ", msg: "", ok: true},
		{line: "   if not rst_n_i  then ", msg: "", ok: true},
	}

	for i, test := range tests {
		msg, ok := checkNegativeResetIfCondition([]byte(test.line))
		if msg != test.msg || ok != test.ok {
			t.Errorf("[%d]: got (%v, %v); want (%v, %v)", i, msg, ok, test.msg, test.ok)
		}
	}
}
