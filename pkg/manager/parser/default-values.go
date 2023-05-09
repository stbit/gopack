package parser

var zeroValue = []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune", "float32", "float64", "complex64"}

func getDefaultValue(t string) string {
	for _, v := range zeroValue {
		if v == t {
			return "0"
		}
	}

	switch t {
	case "string":
		return "\"\""
	case "bool":
		return "false"
	default:
		return t + "{}"
	}
}
