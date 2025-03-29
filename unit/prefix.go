package unit

import "maps"

// Prefix represents a unit prefix.
type Prefix struct {
	Name   string // Name is the name of the unit prefix, e.g. "kilo".
	Symbol string // Symbol is the short form of the unit symbol, e.g. "k".
}

var decimalPrefixes = map[float64]Prefix{
	Kilo:   {"kilo", "k"},
	Mega:   {"mega", "M"},
	Giga:   {"giga", "G"},
	Tera:   {"tera", "T"},
	Peta:   {"peta", "P"},
	Exa:    {"exa", "E"},
	Zetta:  {"zetta", "Z"},
	Yotta:  {"yotta", "Y"},
	Ronna:  {"ronna", "R"},
	Quetta: {"quetta", "Q"},
}

var binaryPrefixes = map[float64]Prefix{
	Kibi:  {"kibi", "Ki"},
	Mebi:  {"mebi", "Mi"},
	Gibi:  {"gibi", "Gi"},
	Tebi:  {"tebi", "Ti"},
	Pebi:  {"pebi", "Pi"},
	Exbi:  {"exbi", "Ei"},
	Zebi:  {"zebi", "Zi"},
	Yobi:  {"yobi", "Yi"},
	Robi:  {"robi", "Ri"},
	Quebi: {"quebi", "Qi"},
}

var prefixes = func() map[float64]Prefix {
	combined := make(map[float64]Prefix)
	maps.Copy(combined, decimalPrefixes)
	maps.Copy(combined, binaryPrefixes)
	return combined
}()

// PrefixFor returns the prefix for a given value.
func PrefixFor(value float64) Prefix {
	if v, found := prefixes[value]; found {
		return v
	}
	return Prefix{}
}

// prefixFromMap safely retrieves a prefix from a prefix map.
// It returns the matching prefix for the given value or an empty prefix if not
// found.
func prefixFromMap(value float64, prefixMap map[float64]Prefix) Prefix {
	if prefix, found := prefixMap[value]; found {
		return prefix
	}

	return Prefix{}
}
