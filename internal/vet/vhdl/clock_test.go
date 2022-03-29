package vhdl

import (
	"testing"
)

func TestCheckClockPortMapping(t *testing.T) {
	var tests = []struct {
		line string
		msg  string
		ok   bool
	}{
		// Invalid mappings
		{line: "clk_10=>clk_20", msg: "clock frequency mismatch", ok: false},
		{line: "clk_40_i => clk_160)", msg: "clock frequency mismatch", ok: false},
		{line: "clk70 => clk_80", msg: "clock frequency mismatch", ok: false},
		{line: "clk_70 => clk80_i", msg: "clock frequency mismatch", ok: false},
		{line: "clk70 => clock120", msg: "clock frequency mismatch", ok: false},
		{line: "clock70_i => clk120_i,", msg: "clock frequency mismatch", ok: false},
		// Valid mappings
		{line: "clk70_i => clk70,", msg: "", ok: true},
		{line: "clk25 => clk_25,", msg: "", ok: true},
		{line: "clock125 => clk_125,", msg: "", ok: true},
		{line: "clock_40_i=>clock40,", msg: "", ok: true},
	}

	for i, test := range tests {
		msg, ok := checkClockPortMapping([]byte(test.line))
		if msg != test.msg || ok != test.ok {
			t.Errorf("[%d]: got (%v, %v); want (%v, %v)", i, msg, ok, test.msg, test.ok)
		}
	}
}
