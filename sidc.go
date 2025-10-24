package milstd2525c

import (
	"fmt"
	"strings"
)

// SIDC represents a Symbol Identification Code (15 characters)
// Format: [Coding Scheme][Std Identity][Battle Dimension][Status][Function ID (6 chars)][Modifier (5 chars)]
type SIDC struct {
	CodingScheme     CodingScheme     // Position 1
	StandardIdentity StandardIdentity // Position 2
	BattleDimension  BattleDimension  // Position 3
	Status           Status           // Position 4
	FunctionID       string           // Positions 5-10 (6 characters)
	Modifier         string           // Positions 11-15 (5 characters)
}

// CodingScheme represents the symbology set (Position 1)
type CodingScheme string

const (
	CodingSchemeWarfighting CodingScheme = "S" // Warfighting/Units/Equipment
	CodingSchemeOperations  CodingScheme = "G" // Military Operations
	CodingSchemeMETOC       CodingScheme = "W" // Meteorological and Oceanographic
	CodingSchemeSIGINT      CodingScheme = "I" // Signals Intelligence
	CodingSchemeEmergency   CodingScheme = "E" // Emergency Management
)

// StandardIdentity represents the threat posed by the object (Position 2)
type StandardIdentity string

const (
	StandardIdentityPending StandardIdentity = "P" // Pending
	StandardIdentityUnknown StandardIdentity = "U" // Unknown
	StandardIdentityFriend  StandardIdentity = "F" // Friend
	StandardIdentityNeutral StandardIdentity = "N" // Neutral
	StandardIdentityHostile StandardIdentity = "H" // Hostile
	StandardIdentitySuspect StandardIdentity = "S" // Suspect
	StandardIdentityJoker   StandardIdentity = "J" // Joker (Exercise Hostile)
)

// BattleDimension represents the primary mission area (Position 3)
type BattleDimension string

const (
	BattleDimensionUnknown    BattleDimension = "Z" // Unknown
	BattleDimensionSpace      BattleDimension = "P" // Space
	BattleDimensionAir        BattleDimension = "A" // Air
	BattleDimensionGround     BattleDimension = "G" // Ground
	BattleDimensionSeaSurface BattleDimension = "S" // Sea Surface
	BattleDimensionSubsurface BattleDimension = "U" // Subsurface
	BattleDimensionSOF        BattleDimension = "F" // Special Operations Forces
)

// Status represents whether the object is present or planned (Position 4)
type Status string

const (
	StatusPresent   Status = "P" // Present/Existing
	StatusPlanned   Status = "A" // Anticipated/Planned
	StatusDamaged   Status = "D" // Damaged
	StatusDestroyed Status = "X" // Destroyed
)

// ParseSIDC parses a 15-character SIDC string
func ParseSIDC(code string) (*SIDC, error) {
	code = strings.TrimSpace(code)

	if len(code) != 15 {
		return nil, fmt.Errorf("SIDC must be exactly 15 characters, got %d", len(code))
	}

	return &SIDC{
		CodingScheme:     CodingScheme(code[0:1]),
		StandardIdentity: StandardIdentity(code[1:2]),
		BattleDimension:  BattleDimension(code[2:3]),
		Status:           Status(code[3:4]),
		FunctionID:       code[4:10],
		Modifier:         code[10:15],
	}, nil
}

// String returns the 15-character SIDC code
func (s *SIDC) String() string {
	return fmt.Sprintf("%s%s%s%s%s%s",
		s.CodingScheme,
		s.StandardIdentity,
		s.BattleDimension,
		s.Status,
		s.FunctionID,
		s.Modifier,
	)
}
