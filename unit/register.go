package unit

import (
	"errors"
	"slices"
	"sync"
)

var (
	// Next available unit ID, used for custom units.
	nextUnitID = unitMaxBuiltin

	unitRegistryMutex  sync.RWMutex
	formatSystemsMutex sync.RWMutex

	// Registry used for storing format systems.
	formatSystems = make(map[string]*FormatSystem)

	// Reverse lookup maps used for performant lookups (Protected by unitRegistryMutex)
	unitSymbolLookup   = make(map[string]Unit)
	unitSingularLookup = make(map[string]Unit)
	unitPluralLookup   = make(map[string]Unit)

	// Combined prefix lookup maps used for performant lookups (Protected by formatSystemsMutex)
	// Maps prefix Symbol/Name string to its scale factor
	prefixLookup = make(map[string]float64)

	// Keep track of known prefix strings sorted by length descending for matching
	sortedPrefixKeys = []string{}
)

// Register creates a new custom unit with the provided descriptor.
// It returns a Unit value that can be used with the formatting functions.
func Register(descriptor Descriptor) (Unit, error) {
	unitRegistryMutex.Lock()
	defer unitRegistryMutex.Unlock()

	// Validate the descriptor
	if descriptor.Symbol == "" && descriptor.Singular == "" && descriptor.Plural == "" {
		return 0, errors.New("unit: RegisterUnit requires non-empty descriptor fields")
	}

	// Allocate a new unit ID
	id := nextUnitID
	nextUnitID++

	// Register the unit
	unitRegistry[id] = descriptor

	// Update lookup maps
	rebuildUnitLookups()

	return id, nil
}

// MustRegister is like Register but panics on error
func MustRegister(descriptor Descriptor) Unit {
	id, err := Register(descriptor)
	if err != nil {
		panic(err)
	}
	return id
}

// FormatSystem represents a system of units with associated boundaries
type FormatSystem struct {
	// Name is the unique identifier for the format system, used for registration
	// and retrieval (e.g., via GetFormatSystem). Examples: "decimal", "binary", "time".
	Name string
	// Boundaries define the scaling thresholds for applying prefixes. It's a slice
	// of float64 values, typically sorted in descending order (e.g., [Mega, Kilo]
	// or [Mebi, Kibi]). A value is divided by the largest boundary it is
	// greater than or equal to during formatting.
	Boundaries []float64
	// Prefixes maps a boundary value (from the Boundaries slice) to its
	// corresponding Prefix definition (Name and Symbol). This map provides the
	// textual representation for each scaling level defined by the boundaries.
	// Example: {1000.0: {Name: "kilo", Symbol: "k"}, 1024.0: {Name: "kibi", Symbol: "Ki"}}
	Prefixes map[float64]Prefix
}

// RegisterFormatSystem registers a new formatting system
func RegisterFormatSystem(name string, boundaries []float64, prefixes map[float64]Prefix) (*FormatSystem, error) {
	formatSystemsMutex.Lock()
	defer formatSystemsMutex.Unlock()

	// Validate input parameters
	if name == "" || prefixes == nil || len(boundaries) != len(prefixes) || len(boundaries) == 0 {
		return nil, errors.New("unit: empty format system name or nil prefixes")
	}

	system := &FormatSystem{
		Name:       name,
		Boundaries: boundaries,
		Prefixes:   prefixes,
	}

	// Register with global registry
	formatSystems[name] = system

	// Update lookup maps
	rebuildPrefixLookups()

	return system, nil
}

// MustRegisterFormatSystem is like RegisterFormatSystem but panics on error
func MustRegisterFormatSystem(name string, boundaries []float64, prefixes map[float64]Prefix) *FormatSystem {
	system, err := RegisterFormatSystem(name, boundaries, prefixes)
	if err != nil {
		panic(err)
	}
	return system
}

// GetFormatSystem retrieves a registered format system by name
func GetFormatSystem(name string) (*FormatSystem, bool) {
	formatSystemsMutex.RLock()
	defer formatSystemsMutex.RUnlock()

	system, found := formatSystems[name]
	return system, found
}

// BuiltinFormatSystems returns a map of all built-in format systems
func BuiltinFormatSystems() map[string]*FormatSystem {
	return map[string]*FormatSystem{
		"decimal": {
			Name:       "decimal",
			Boundaries: SIDecimalBoundaries,
			Prefixes:   decimalPrefixes,
		},
		"binary": {
			Name:       "binary",
			Boundaries: IECBinaryBoundaries,
			Prefixes:   binaryPrefixes,
		},
	}
}

func rebuildUnitLookups() {
	unitSymbolLookup = make(map[string]Unit)
	unitSingularLookup = make(map[string]Unit)
	unitPluralLookup = make(map[string]Unit)
	for u, desc := range unitRegistry {
		if desc.Symbol != "" {
			unitSymbolLookup[desc.Symbol] = u
		}
		if desc.Singular != "" {
			unitSingularLookup[desc.Singular] = u
		}
		if desc.Plural != "" {
			unitPluralLookup[desc.Plural] = u
		}
	}
}

func rebuildPrefixLookups() {
	prefixLookup = make(map[string]float64)
	// Add built-in SI prefixes
	for scale, prefix := range decimalPrefixes {
		if prefix.Symbol != "" {
			prefixLookup[prefix.Symbol] = scale // Adds 'k', 'M', 'G', etc.
		}
		// Explicitly add uppercase 'K' as an alias for Kilo
		if scale == Kilo {
			prefixLookup["K"] = scale // Add 'K' -> 1000
		}
		if prefix.Name != "" {
			prefixLookup[prefix.Name] = scale // Adds 'kilo', 'mega', etc.
		}
	}
	// Add built-in IEC prefixes
	for scale, prefix := range binaryPrefixes {
		if prefix.Symbol != "" {
			prefixLookup[prefix.Symbol] = scale // Adds 'Ki', 'Mi', etc.
		}
		if prefix.Name != "" {
			prefixLookup[prefix.Name] = scale // Adds 'kibi', 'mebi', etc.
		}
	}
	// Add prefixes from custom registered systems
	for _, system := range formatSystems {
		for scale, prefix := range system.Prefixes {
			// Avoid overwriting potentially more common built-ins? Or let last registration win?
			// Current approach: Last registration wins within custom systems.
			if prefix.Symbol != "" {
				prefixLookup[prefix.Symbol] = scale
			}
			if prefix.Name != "" {
				prefixLookup[prefix.Name] = scale
			}
		}
	}

	// Rebuild sorted keys for matching
	sortedPrefixKeys = make([]string, 0, len(prefixLookup))
	for k := range prefixLookup {
		sortedPrefixKeys = append(sortedPrefixKeys, k)
	}
	// Sort descending by length - crucial for longest match first!
	slices.SortFunc(sortedPrefixKeys, func(a, b string) int {
		return len(b) - len(a) // Descending length
	})
}

func init() {
	// Populate unit lookups from built-ins
	rebuildUnitLookups()

	// Populate prefix lookups from built-ins
	rebuildPrefixLookups()
}
