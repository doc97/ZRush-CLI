package zsteps

import (
	"fmt"

	"github.com/doc97/zrush/zdata"
	"github.com/doc97/zrush/zutil"
)

// StepAttack performs the attack step
func StepAttack(player *zdata.PlayerData, playerCount int, units map[string]zdata.UnitData) (*zdata.AttackData, error) {
	if player.Units[zdata.Zergling] == 0 && player.Units[zdata.Hydralisk] == 0 && player.Units[zdata.Mutalisk] == 0 {
		return nil, nil
	}

	subStepPrintAttackUnits(player, units)

	attackUnits, err := subStepSelectAttackUnits(player)
	if err != nil {
		return nil, fmt.Errorf("failed to select units: %v", err)
	}

	targetPlayer := subStepSelectAttackTarget(playerCount, player.ID)
	attack := zdata.AttackData{
		AttackerID: player.ID,
		DefenderID: targetPlayer,
		Units:      attackUnits,
	}

	fmt.Println()
	return &attack, nil
}

func subStepPrintAttackUnits(player *zdata.PlayerData, units map[string]zdata.UnitData) {
	fmt.Println("Available offensive units:")

	nonAttackUnits := map[string]struct{}{zdata.Drone: {}, zdata.SporeCrawler: {}}
	for key, count := range player.Units {
		if _, contains := nonAttackUnits[key]; contains {
			continue
		}

		if count == 1 {
			fmt.Printf("\t%d %s\n", count, units[key].Name)
		} else if count > 1 {
			fmt.Printf("\t%d %ss\n", count, units[key].Name)
		}
	}
}

func subStepSelectAttackUnits(player *zdata.PlayerData) (map[string]int, error) {
	fmt.Printf("\nAttack with (format: 'zzh')? ")
	units, _, err := zutil.ReadUnitString()
	if err != nil {
		return nil, err
	}
	for name := range units {
		units[name] = zutil.Min(units[name], player.Units[name])
	}
	return units, nil
}

func subStepSelectAttackTarget(playerCount, selfID int) int {
	targetPlayer := -1
	for {
		fmt.Printf("Attack player [1-%d]? ", playerCount)
		var target int
		if _, err := fmt.Scan(&target); err != nil {
			continue
		}
		if target == selfID || target < 1 || target > playerCount {
			continue
		}
		targetPlayer = target
		break
	}
	return targetPlayer
}
