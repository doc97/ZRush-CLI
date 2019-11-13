package zutil

import (
	"fmt"
	"strings"

	"github.com/doc97/zrush/zdata"
)

// ReadUnitString prompts the user for a unit string and returns a map of units and
// their counts, the inputted string and a possible error.
func ReadUnitString() (map[string]int, string, error) {
	var unitStr string
	if _, err := fmt.Scan(&unitStr); err != nil {
		return nil, "", fmt.Errorf("could not read input: %v", err)
	}

	units := map[string]int{
		zdata.Drone:        0,
		zdata.Zergling:     0,
		zdata.Hydralisk:    0,
		zdata.Mutalisk:     0,
		zdata.SporeCrawler: 0,
	}
	for _, unitRune := range strings.ToLower(unitStr) {
		switch unitRune {
		case 'd':
			units[zdata.Drone]++
		case 'z':
			units[zdata.Zergling]++
		case 'h':
			units[zdata.Hydralisk]++
		case 'm':
			units[zdata.Mutalisk]++
		case 's':
			units[zdata.SporeCrawler]++
		}
	}

	return units, unitStr, nil
}
