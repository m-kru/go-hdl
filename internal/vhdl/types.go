package vhdl

func IsSingleBitStdType(typ string) bool {
	singleBitStdTypes := map[string]bool{
		"bit": true, "boolean": true, "std_logic": true, "std_ulogic": true,
	}
	if _, ok := singleBitStdTypes[typ]; ok {
		return true
	}
	return false
}
