package milstd2525c

import (
	"testing"
)

func TestParseSIDC(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *SIDC
		wantErr  bool
	}{
		{
			name:  "Friendly Ground Infantry",
			input: "SFGPUCII-------",
			expected: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPresent,
				FunctionID:       "UCII--",
				Modifier:         "-----",
			},
			wantErr: false,
		},
		{
			name:  "Hostile Air Fighter",
			input: "SHAPMFF--------",
			expected: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityHostile,
				BattleDimension:  BattleDimensionAir,
				Status:           StatusPresent,
				FunctionID:       "MFF---",
				Modifier:         "-----",
			},
			wantErr: false,
		},
		{
			name:  "Unknown Subsurface Submarine",
			input: "SUUPSU---------",
			expected: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityUnknown,
				BattleDimension:  BattleDimensionSubsurface,
				Status:           StatusPresent,
				FunctionID:       "SU----",
				Modifier:         "-----",
			},
			wantErr: false,
		},
		{
			name:  "Neutral Sea Surface Cruiser",
			input: "SNSPCL---------",
			expected: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityNeutral,
				BattleDimension:  BattleDimensionSeaSurface,
				Status:           StatusPresent,
				FunctionID:       "CL----",
				Modifier:         "-----",
			},
			wantErr: false,
		},
		{
			name:  "Planned Friendly Ground Armor",
			input: "SFGAUCIC-------",
			expected: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPlanned,
				FunctionID:       "UCIC--",
				Modifier:         "-----",
			},
			wantErr: false,
		},
		{
			name:    "Invalid - Too Short",
			input:   "SFGP",
			wantErr: true,
		},
		{
			name:    "Invalid - Too Long",
			input:   "SFGPUCII----------",
			wantErr: true,
		},
		{
			name:    "Invalid - Empty",
			input:   "",
			wantErr: true,
		},
		{
			name:  "General Ground Unit",
			input: "SFGPU----------",
			expected: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPresent,
				FunctionID:       "U-----",
				Modifier:         "-----",
			},
			wantErr: false,
		},
		{
			name:  "Unknown Fixed Wing Drone",
			input: "SUAPMFQ--------",
			expected: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityUnknown,
				BattleDimension:  BattleDimensionAir,
				Status:           StatusPresent,
				FunctionID:       "MFQ---",
				Modifier:         "-----",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSIDC(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSIDC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr {
				return
			}
			
			if got.CodingScheme != tt.expected.CodingScheme {
				t.Errorf("CodingScheme = %v, want %v", got.CodingScheme, tt.expected.CodingScheme)
			}
			if got.StandardIdentity != tt.expected.StandardIdentity {
				t.Errorf("StandardIdentity = %v, want %v", got.StandardIdentity, tt.expected.StandardIdentity)
			}
			if got.BattleDimension != tt.expected.BattleDimension {
				t.Errorf("BattleDimension = %v, want %v", got.BattleDimension, tt.expected.BattleDimension)
			}
			if got.Status != tt.expected.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.expected.Status)
			}
			if got.FunctionID != tt.expected.FunctionID {
				t.Errorf("FunctionID = %v, want %v", got.FunctionID, tt.expected.FunctionID)
			}
			if got.Modifier != tt.expected.Modifier {
				t.Errorf("Modifier = %v, want %v", got.Modifier, tt.expected.Modifier)
			}
		})
	}
}

func TestSIDCString(t *testing.T) {
	tests := []struct {
		name     string
		sidc     *SIDC
		expected string
	}{
		{
			name: "Friendly Ground Infantry",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPresent,
				FunctionID:       string(FunctionGrdtrkUntCbtInfInffv),
				Modifier:         "-----",
			},
			expected: "SFGPUCII-------",
		},
		{
			name: "Hostile Air Fighter",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityHostile,
				BattleDimension:  BattleDimensionAir,
				Status:           StatusPresent,
				FunctionID:       string(FunctionAirtrkMilFixdFtr),
				Modifier:         "-----",
			},
			expected: "SHAPMFF--------",
		},
		{
			name: "Unknown UAV with Modifiers",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityUnknown,
				BattleDimension:  BattleDimensionAir,
				Status:           StatusPresent,
				FunctionID:       string(FunctionAirtrkMilFixdUty),
				Modifier:         "12345",
			},
			expected: "SUAPMFU---12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sidc.String()
			if got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	testCases := []string{
		"SFGPUCII-------", // Friendly Ground Infantry
		"SHAPMFF--------", // Hostile Air Fighter
		"SUUPSU---------", // Unknown Subsurface Submarine
		"SNSPCL---------", // Neutral Sea Surface Cruiser
		"SFGPU----------", // General Ground Unit
		"SUAPMFQ--------", // Unknown Fixed Wing Drone
		"EUUP------12345", // Emergency Management
		"SGAP------12345", // Operations with modifiers
		"WPAP------ABCDE", // METOC with modifiers
	}

	for _, original := range testCases {
		t.Run(original, func(t *testing.T) {
			// Parse the SIDC
			parsed, err := ParseSIDC(original)
			if err != nil {
				t.Fatalf("Failed to parse SIDC %s: %v", original, err)
			}

			// Convert back to string
			result := parsed.String()
			
			// Should match the original
			if result != original {
				t.Errorf("Round trip failed: original=%s, result=%s", original, result)
			}
		})
	}
}

func BenchmarkParseSIDC(b *testing.B) {
	code := "SFGPUCII-------"
	for i := 0; i < b.N; i++ {
		_, _ = ParseSIDC(code)
	}
}

func BenchmarkSIDCString(b *testing.B) {
	sidc := &SIDC{
		CodingScheme:     CodingSchemeWarfighting,
		StandardIdentity: StandardIdentityFriend,
		BattleDimension:  BattleDimensionGround,
		Status:           StatusPresent,
		FunctionID:       string(FunctionGrdtrkUntCbtInfInffv),
		Modifier:         "-----",
	}
	
	for i := 0; i < b.N; i++ {
		_ = sidc.String()
	}
}