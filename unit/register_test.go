package unit_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/require"
	"github.com/neticdk/go-stdlib/unit"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  unit.Descriptor
		expectError bool
	}{
		{
			name: "valid unit",
			descriptor: unit.Descriptor{
				Symbol:   "W",
				Singular: "watt",
				Plural:   "watts",
			},
			expectError: false,
		},
		{
			name: "empty symbol",
			descriptor: unit.Descriptor{
				Symbol:   "",
				Singular: "watt",
				Plural:   "watts",
			},
			expectError: false,
		},
		{
			name: "empty singular",
			descriptor: unit.Descriptor{
				Symbol:   "W",
				Singular: "",
				Plural:   "watts",
			},
			expectError: false,
		},
		{
			name: "empty plural",
			descriptor: unit.Descriptor{
				Symbol:   "W",
				Singular: "watt",
				Plural:   "",
			},
			expectError: false,
		},
	}

	// Store initial registered units to compare later
	initialUnits := unit.List()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, err := unit.Register(tt.descriptor)

			if tt.expectError {
				assert.Error(t, err, "Register() should have returned an error")
			} else {
				assert.NoError(t, err, "Register() should not have returned an error")

				// Check the descriptor matches what we registered
				desc := unit.Describe(kind)
				assert.Equal(t, tt.descriptor.Symbol, desc.Symbol, "Symbol should match")
				assert.Equal(t, tt.descriptor.Singular, desc.Singular, "Singular name should match")
				assert.Equal(t, tt.descriptor.Plural, desc.Plural, "Plural name should match")
			}
		})
	}

	// Check that new units were added to registry
	afterUnits := unit.List()
	assert.Greater(t, len(afterUnits), len(initialUnits), "Register() should have added new units")
}

func TestMustRegister(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  unit.Descriptor
		shouldPanic bool
	}{
		{
			name: "valid unit",
			descriptor: unit.Descriptor{
				Symbol:   "A",
				Singular: "ampere",
				Plural:   "amperes",
			},
			shouldPanic: false,
		},
		{
			name: "invalid unit",
			descriptor: unit.Descriptor{
				Symbol:   "",
				Singular: "",
				Plural:   "",
			},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					unit.MustRegister(tt.descriptor)
				}, "MustRegister() should have panicked with invalid descriptor")
			} else {
				assert.NotPanics(t, func() {
					kind := unit.MustRegister(tt.descriptor)

					// Verify the unit was registered correctly
					desc := unit.Describe(kind)
					assert.Equal(t, tt.descriptor, desc, "Descriptor should match what was registered")
				}, "MustRegister() should not panic with valid descriptor")
			}
		})
	}
}

func TestUniqueUnitIDs(t *testing.T) {
	// Register multiple units and ensure they get different IDs
	unit1, err1 := unit.Register(unit.Descriptor{Symbol: "U1", Singular: "unit1", Plural: "units1"})
	unit2, err2 := unit.Register(unit.Descriptor{Symbol: "U2", Singular: "unit2", Plural: "units2"})
	unit3, err3 := unit.Register(unit.Descriptor{Symbol: "U3", Singular: "unit3", Plural: "units3"})

	// Ensure registrations succeeded
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)

	// Check for uniqueness
	assert.NotEqual(t, unit1, unit2, "First and second units should have different IDs")
	assert.NotEqual(t, unit1, unit3, "First and third units should have different IDs")
	assert.NotEqual(t, unit2, unit3, "Second and third units should have different IDs")
}

func TestRegisteredUnitsInList(t *testing.T) {
	// Register a new unit
	desc := unit.Descriptor{Symbol: "J", Singular: "joule", Plural: "joules"}
	kind, err := unit.Register(desc)
	require.NoError(t, err, "Registration should succeed")

	// Check if the unit appears in List()
	unitList := unit.List()

	found := false
	for _, info := range unitList {
		if info.Unit == kind {
			assert.Equal(t, desc.Symbol, info.Symbol, "Symbol should match")
			assert.Equal(t, desc.Singular, info.Singular, "Singular should match")
			assert.Equal(t, desc.Plural, info.Plural, "Plural should match")
			found = true
			break
		}
	}

	assert.True(t, found, "Registered unit should be found in List()")
}
