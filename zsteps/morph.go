package zsteps

import (
	"fmt"

	"github.com/doc97/zrush/zdata"
	"github.com/doc97/zrush/zutil"
)

// StepMorph performs the morph step
func StepMorph(player *zdata.PlayerData, units map[string]zdata.UnitData) error {
	for {
		// Exit condition
		if !subStepPrintAffordableUnits(player, units) {
			fmt.Printf("Morphing units...\n\n")
			return nil
		}

		fmt.Printf("\nMorph a unit, 'x' to skip [d,z,h,m,s]? ")
		morphUnits, str, err := zutil.ReadUnitString()
		if err != nil {
			return err
		}
		if str == "x" {
			return nil
		}

		costMinerals, costGas := 0, 0
		for name, count := range morphUnits {
			costMinerals += count * units[name].Minerals
			costGas += count * units[name].Gas
		}

		if costMinerals > player.Minerals || costGas > player.Gas {
			fmt.Println("Insufficient funds!")
			continue
		}

		for name, count := range morphUnits {
			player.Units[name] += count
		}

		player.Minerals -= costMinerals
		player.Gas -= costGas
	}
}

func subStepPrintAffordableUnits(player *zdata.PlayerData, units map[string]zdata.UnitData) bool {
	affordableUnits := make(map[string]zdata.UnitData)
	for name, unit := range units {
		if player.Minerals >= unit.Minerals && player.Gas >= unit.Gas {
			affordableUnits[name] = unit
		}
	}

	if len(affordableUnits) == 0 {
		return false
	}

	fmt.Printf("Minerals: %d\nVespene Gas: %d\n\n", player.Minerals, player.Gas)
	fmt.Println("You can morph the following units:")
	for _, unit := range affordableUnits {
		fmt.Printf("%s (%d minerals, %d gas)\n", unit.Name, unit.Minerals, unit.Gas)
	}

	return true
}
