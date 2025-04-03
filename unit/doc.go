// Package unit provides utilities for handling units including human-readable
// unit formatting, parsing, and registration of custom units and formatting systems.
//
// Unit Systems:
//
// This package supports two built-in unit systems:
//
// 1. IEC Binary Units (base-2): Used by Binary() option and FormatBinary functions.
//   - Units increase by powers of 1024 (2^10)
//   - Prefixes include Ki, Mi, Gi, Ti, Pi, etc.
//   - Often used for memory and storage in computing contexts
//   - Example: 1024 bytes = 1 KiB (kibibyte)
//
// 2. SI Decimal Units (base-10): Used by Decimal() option and FormatDecimal functions.
//   - Units increase by powers of 1000 (10^3)
//   - Prefixes include k, M, G, T, P, etc.
//   - Often used for data transmission rates and disk marketing
//   - Example: 1000 bytes = 1 kB (kilobyte)
//
// Supported Units:
//
// The package natively supports the following unit types:
//   - Byte: Used for data storage (B)
//   - Bit: Used for data transmission (b)
//   - Hertz: Used for frequencies (Hz)
//
// Formatting Functions:
//
// This package provides multiple approaches to formatting values with units:
//
// 1. Option-Based Approach (Most Flexible):
//   - Format: Returns a Value struct that provides both the scaled value
//     and full formatting information. Configure with option functions like
//     Binary(), Decimal(), WithUnit(), WithPrecision(), etc.
//
// Examples:
//
//	// Basic usage
//	fmt.Println(Format(1024, Binary(), WithUnit(Byte))) // Output: 1 KiB
//
//	// Multiple options
//	fmt.Println(Format(1500,
//	    Decimal(),
//	    WithUnit(Bit),
//	    WithPrecision(2),
//	    WithPlural())) // Output: 1.50 kilobits
//
// 2. Simplified Convenience Functions:
//   - FormatBinary/FormatDecimal: Format values without a unit type using standard
//     binary or decimal systems.
//   - FormatBinaryUnit/FormatDecimalUnit: Format with specific built-in unit types
//     (Byte, Bit, Hertz) using standard binary or decimal systems.
//
// These convenience functions internally use Format with specific option
// combinations.
//
// Returned Value:
//
// The Format function returns a Value struct which includes:
//   - Scaled: The scaled numeric value.
//   - Prefix: Information about the unit prefix (e.g., "kilo", "Ki").
//   - Unit: A Descriptor containing information about the unit (e.g., "byte", "bit").
//
// The Value struct implements fmt.Stringer, so it can be used directly in
// string contexts like fmt.Printf or fmt.Println.
//
// Custom Units:
//
// Beyond the built-in units (Byte, Bit, Hertz), you can register your own custom units.
// This allows the formatting functions to work with units specific to your domain.
//
//   - Register / MustRegister: Use these functions to define a new unit.
//     You provide a Descriptor containing the unit's Symbol, Singular name,
//     and Plural name. MustRegister panics on error, while Register returns an error.
//   - Descriptor: This struct holds the textual representations of a unit.
//     It also includes the Unit identifier itself.
//   - List / Names / Describe / Unit.Info: Functions to query registered units.
//
// Example:
//
//	// Register a custom "pixel" unit
//	UnitPixel := unit.MustRegister(unit.Descriptor{
//	    Symbol:   "px",
//	    Singular: "pixel",
//	    Plural:   "pixels",
//	})
//
//	// Use the custom unit in formatting
//	fmt.Println(unit.Format(2048, unit.Decimal(), unit.WithUnit(UnitPixel))) // Output: 2 kpx
//	fmt.Println(unit.Format(1, unit.Decimal(), unit.WithUnit(UnitPixel)))    // Output: 1 px
//
// Custom Formatting Systems:
//
// If the standard SI (base-10) or IEC (base-2) scaling doesn't fit your needs,
// you can define completely custom formatting systems with different boundaries
// and prefix names/symbols.
//
//   - FormatSystem: Represents a custom system with its name, scaling boundaries,
//     and prefix map.
//   - RegisterFormatSystem: Registers a new FormatSystem globally by name.
//   - GetFormatSystem: Retrieves a registered system by name.
//   - WithSystem / WithSystemByName: FormatOptions to apply a custom system
//     during formatting with Format.
//
// Example:
//
//	// Define a simple time formatting system (seconds, minutes, hours)
//	timeSystem := unit.RegisterFormatSystem(
//	    "time",
//	    []float64{3600, 60}, // Boundaries (hour, minute)
//	    map[float64]unit.Prefix{
//	        3600: {Name: "hour", Symbol: "h"},
//	        60:   {Name: "minute", Symbol: "m"},
//	    },
//	)
//	UnitSecond := unit.MustRegister(unit.Descriptor{Symbol: "", Singular: "", Plural: "s"})
//
//	// Format a duration using the custom system
//	duration := 3725.0 // seconds
//	fmt.Println(unit.Format(duration, unit.WithSystem(timeSystem), unit.WithUnit(UnitSecond), unit.WithPrecision(1))) // Output: 1.0 h
//
// Parsing:
//
// The Parse function attempts to interpret strings containing numbers and
// standard units (e.g., "1.5 KB", "1024", "500 MHz"). It currently only
// supports built-in SI/IEC prefixes and the built-in units (Byte, Bit, Hertz).
// Custom registered units and custom format systems are **not** supported by
// the Parse function in this version.
//
// When to Use Each System:
//
//   - Use binary (IEC) units when working with computer memory, storage capacities,
//     or when communicating with technical audiences.
//   - Use decimal (SI) units when dealing with data transfer rates, disk marketing,
//     or when communicating with general audiences.
//
// Choose the appropriate formatting approach and extensibility features based on
// your specific use case and audience expectations.
package unit
