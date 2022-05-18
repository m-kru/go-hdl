package vhdl

import (
	"bufio"
	"bytes"
	"testing"
)

func TestEnumScanning(t *testing.T) {
	var tests = []struct {
		code   string
		enum   enum
		width  uint
		values []string
	}{
		{
			code:   `type t_state is (ONE, TWO);`,
			enum:   enum{name: "t_state", values: []string{"ONE", "TWO"}},
			width:  1,
			values: []string{"ONE", "TWO"},
		},
		{
			code: `type t_state is (
                      ONE, TWO ) ;`,
			enum:   enum{name: "t_state", values: []string{"ONE", "TWO"}},
			width:  1,
			values: []string{"ONE", "TWO"},
		},
		{
			code: `type t_state is ( ONE
                     TWO, THREE
                   );  `,
			enum:   enum{name: "t_state", values: []string{"ONE", "TWO"}},
			width:  2,
			values: []string{"ONE", "TWO", "THREE"},
		},
		{
			code: `type t_state is (
                      ONE,
                      TWO,
                      THREE
                   );  `,
			enum:   enum{name: "t_state", values: []string{"ONE", "TWO"}},
			width:  2,
			values: []string{"ONE", "TWO", "THREE"},
		},
	}

	for i, test := range tests {
		sCtx := scanContext{scanner: bufio.NewScanner(bytes.NewReader([]byte(test.code)))}
		sCtx.proceed()
		enum, err := scanEnumTypeDeclaration(&sCtx, test.enum.name)
		if err != nil {
			t.Errorf("[%d]: %v", i, err)
		}
		if len(enum.values) != len(test.values) {
			t.Errorf("[%d]: invalid values length, got %d, want %d", i, len(enum.values), len(test.values))
		}
		for j, v := range enum.values {
			if v != test.values[j] {
				t.Errorf("[%d]: invalid value %s, want %s", i, v, test.values[j])
			}
		}
		if enum.Width() != test.width {
			t.Errorf("[%d]: invalid width %d, want %d", i, enum.Width(), test.width)
		}
	}
}
