package milstd2525c

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
)

// MilSymbolResult represents the result from milsymbol symbol construction and validation
type MilSymbolResult struct {
	SIDC                    string `json:"sidc"`
	IsValid                 bool   `json:"isValid"`
	Error                   string `json:"error,omitempty"`
	HasVisualRepresentation bool   `json:"hasVisualRepresentation"`
	SymbolSet               string `json:"symbolSet,omitempty"`
	Entity                  string `json:"entity,omitempty"`
	CanRender               bool   `json:"canRender"`
	SVGLength               int    `json:"svgLength"`
}

// validateWithMilSymbol validates and constructs a SIDC using the milsymbol JavaScript library
func validateWithMilSymbol(sidc string) (*MilSymbolResult, error) {
	script := fmt.Sprintf(`
const { Symbol } = require('milsymbol');

try {
    const symbol = new Symbol('%s');
    
    // Actually construct the symbol and test rendering capabilities
    const svg = symbol.asSVG();
    const isValidSymbol = symbol.isValid();
    
    // Test different rendering methods
    let canRenderSVG = false;
    let canRenderCanvas = false;
    
    try {
        canRenderSVG = typeof symbol.asSVG === 'function' && svg && svg.length > 0;
    } catch (e) {
        canRenderSVG = false;
    }
    
    try {
        canRenderCanvas = typeof symbol.asCanvas === 'function';
        if (canRenderCanvas) {
            const canvas = symbol.asCanvas();
            canRenderCanvas = canvas != null;
        }
    } catch (e) {
        canRenderCanvas = false;
    }
    
    const result = {
        sidc: '%s',
        isValid: isValidSymbol,
        hasVisualRepresentation: canRenderSVG || canRenderCanvas,
        symbolSet: symbol.options ? symbol.options.symbolSet || 'unknown' : 'unknown',
        entity: symbol.options ? symbol.options.entity || 'unknown' : 'unknown',
        canRender: canRenderSVG || canRenderCanvas,
        svgLength: svg ? svg.length : 0
    };
    console.log(JSON.stringify(result));
} catch (error) {
    const result = {
        sidc: '%s',
        isValid: false,
        error: error.message,
        hasVisualRepresentation: false,
        canRender: false,
        svgLength: 0
    };
    console.log(JSON.stringify(result));
}
`, sidc, sidc, sidc)

	cmd := exec.Command("bun", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute milsymbol validation: %w", err)
	}

	var result MilSymbolResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse milsymbol result: %w", err)
	}

	return &result, nil
}

// checkMilSymbolAvailable checks if milsymbol is available for testing
func checkMilSymbolAvailable() error {
	if _, err := exec.LookPath("bun"); err != nil {
		return fmt.Errorf("bun is not available: %w", err)
	}

	cmd := exec.Command("bun", "-e", "const { Symbol } = require('milsymbol'); console.log('OK');")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("milsymbol package not available - run 'bun add milsymbol': %w", err)
	}

	if string(output) != "OK\n" {
		return fmt.Errorf("milsymbol package check failed")
	}

	return nil
}

// TestSIDCValidationAgainstMilSymbol tests essential SIDC cases against milsymbol
func TestSIDCValidationAgainstMilSymbol(t *testing.T) {
	if err := checkMilSymbolAvailable(); err != nil {
		t.Skipf("Skipping milsymbol validation tests: %v", err)
		return
	}

	tests := []struct {
		name        string
		sidc        *SIDC
		expectValid bool
		description string
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
			expectValid: true,
			description: "Standard friendly infantry unit",
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
			expectValid: true,
			description: "Standard hostile fighter aircraft",
		},
		{
			name: "Unknown Submarine",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityUnknown,
				BattleDimension:  BattleDimensionSubsurface,
				Status:           StatusPresent,
				FunctionID:       string(FunctionSbsufSubCnvprn),
				Modifier:         "-----",
			},
			expectValid: true,
			description: "Standard submarine symbol",
		},
		{
			name: "Neutral Cruiser",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityNeutral,
				BattleDimension:  BattleDimensionSeaSurface,
				Status:           StatusPresent,
				FunctionID:       string(FunctionSsufCbttLneCru),
				Modifier:         "-----",
			},
			expectValid: true,
			description: "Standard naval surface ship",
		},
		{
			name: "Friendly UAV",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionAir,
				Status:           StatusPresent,
				FunctionID:       string(FunctionAirtrkMilFixdUty),
				Modifier:         "-----",
			},
			expectValid: true,
			description: "Friendly UAV/drone",
		},
		{
			name: "Suspect UAV",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentitySuspect,
				BattleDimension:  BattleDimensionAir,
				Status:           StatusPresent,
				FunctionID:       string(FunctionAirtrkMilFixdUty),
				Modifier:         "-----",
			},
			expectValid: true,
			description: "Suspect UAV/drone",
		},
		{
			name: "Emergency Management",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeEmergency,
				StandardIdentity: StandardIdentityUnknown,
				BattleDimension:  BattleDimensionUnknown,
				Status:           StatusPresent,
				FunctionID:       "------",
				Modifier:         "12345",
			},
			expectValid: true,
			description: "Emergency Management symbol",
		},
		{
			name: "Operations Symbol",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeOperations,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPresent,
				FunctionID:       "------",
				Modifier:         "12345",
			},
			expectValid: true,
			description: "Military operations with modifiers",
		},
		{
			name: "METOC Symbol",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeMETOC,
				StandardIdentity: StandardIdentityPending,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPresent,
				FunctionID:       "------",
				Modifier:         "ABCDE",
			},
			expectValid: true,
			description: "Weather/oceanographic with modifiers",
		},
		{
			name: "General Ground Unit",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPresent,
				FunctionID:       string(FunctionGrdtrkUntCbt),
				Modifier:         "-----",
			},
			expectValid: true,
			description: "General/unspecified ground unit",
		},
		{
			name: "Unknown Fixed Wing Drone",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityUnknown,
				BattleDimension:  BattleDimensionAir,
				Status:           StatusPresent,
				FunctionID:       FunctionAirtrkMilFixdDrn,
				Modifier:         "-----",
			},
			expectValid: true,
			description: "Unknown fixed wing drone/UAV (MFQ)",
		},
		{
			name: "Invalid Function ID",
			sidc: &SIDC{
				CodingScheme:     CodingSchemeWarfighting,
				StandardIdentity: StandardIdentityFriend,
				BattleDimension:  BattleDimensionGround,
				Status:           StatusPresent,
				FunctionID:       "XXXXXX",
				Modifier:         "-----",
			},
			expectValid: false,
			description: "Invalid function ID should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use our SIDC.String() method to generate the SIDC code
			sidcString := tt.sidc.String()

			result, err := validateWithMilSymbol(sidcString)
			if err != nil {
				t.Fatalf("Failed to validate SIDC %s: %v", sidcString, err)
			}

			if result.IsValid != tt.expectValid {
				t.Errorf("SIDC %s (%s): expected %v, got %v (error: %s)",
					sidcString, tt.description, tt.expectValid, result.IsValid, result.Error)
			}

			// For valid symbols, verify they can actually be constructed and rendered
			if tt.expectValid && result.IsValid {
				if !result.CanRender {
					t.Errorf("SIDC %s: valid symbol should be renderable", sidcString)
				}
				if !result.HasVisualRepresentation {
					t.Logf("SIDC %s: warning - valid symbol has no visual representation", sidcString)
				}
				t.Logf("SIDC %s (%s): valid, renderable=%v, visual=%v, svgLen=%d, symbolSet=%s, entity=%s",
					sidcString, tt.description, result.CanRender, result.HasVisualRepresentation,
					result.SVGLength, result.SymbolSet, result.Entity)
			} else {
				t.Logf("SIDC %s (%s): validation = %v", sidcString, tt.description, result.IsValid)
			}
		})
	}
}

// TestMilSymbolInstallationCheck provides setup instructions if milsymbol is not available
func TestMilSymbolInstallationCheck(t *testing.T) {
	if err := checkMilSymbolAvailable(); err != nil {
		t.Logf(`
To enable milsymbol validation tests: bun add milsymbol
Verify with: bun -e "console.log(require('milsymbol').Symbol)"

Current error: %v
`, err)
	}
}
