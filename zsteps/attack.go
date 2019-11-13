package zsteps

import (
	"fmt"

	"github.com/doc97/zrush/zdata"
	"github.com/doc97/zrush/zutil"
)

// StepAttack performs the attack step
func StepAttack(player *zdata.PlayerData, players []zdata.PlayerData, units map[string]zdata.UnitData) error {
	if player.Units[zdata.Zergling] == 0 && player.Units[zdata.Hydralisk] == 0 && player.Units[zdata.Mutalisk] == 0 {
		return nil
	}

	printAttackUnits(player, units)

	attackUnits, err := selectAttackUnits(player)
	if err != nil {
		return fmt.Errorf("failed to select units: %v", err)
	}

	targetPlayer := selectAttackTarget(players, player.ID)
	attack := zdata.AttackData{
		AttackerID: player.ID,
		Units:      attackUnits,
	}
	targetPlayer.IncomingAttacks = append(targetPlayer.IncomingAttacks, attack)
	fmt.Println()
	return nil
}

func printAttackUnits(player *zdata.PlayerData, units map[string]zdata.UnitData) {
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

func selectAttackUnits(player *zdata.PlayerData) (map[string]int, error) {
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

func selectAttackTarget(players []zdata.PlayerData, selfID int) zdata.PlayerData {
	playerCount := len(players)
	for {
		fmt.Printf("Attack player [1-%d]? ", playerCount)
		var target int
		if _, err := fmt.Scan(&target); err != nil {
			continue
		}
		if target == selfID || target < 1 || target > playerCount {
			continue
		}
		return players[target-1]
	}
}
