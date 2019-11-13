package zdata

// PlayerData holds information about a player state
// like resources and units.
type PlayerData struct {
	ID         int
	BaseHealth int
	Minerals   int
	Gas        int
	Units      map[string]int
}

// AttackData holds information about an incoming attack.
type AttackData struct {
	DefenderID int
	AttackerID int
	Units      map[string]int
}

// UnitData holds information about a unit.
type UnitData struct {
	Name     string
	Attack   int
	Defense  int
	Minerals int
	Gas      int
}

// Drone is the key for Drone units
const Drone = "D"

// Zergling is the key for Zergling units
const Zergling = "z"

// Hydralisk is the key for Hydralisk units
const Hydralisk = "H"

// Mutalisk is the key for Mutalisk units
const Mutalisk = "M"

// SporeCrawler is the key for Spore Crawler units
const SporeCrawler = "SC"
