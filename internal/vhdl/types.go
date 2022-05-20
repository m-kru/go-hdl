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

func IsStdVectorType(typ string) bool {
	stdVectorTypes := map[string]bool{
		"bit_vector":        true,
		"boolean_vector":    true,
		"std_logic_vector":  true,
		"std_ulogic_vector": true,
		"signed":            true,
		"unsigned":          true,
	}
	if _, ok := stdVectorTypes[typ]; ok {
		return true
	}
	return false
}
