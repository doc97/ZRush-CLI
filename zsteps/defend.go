package zsteps

import (
	"fmt"

	"github.com/doc97/zrush/zdata"
)

// StepDefend performs the defend step
func StepDefend(player zdata.PlayerData, players []zdata.PlayerData, units map[string]zdata.UnitData) error {
	if len(player.IncomingAttacks) == 0 {
		return nil
	}

	for attack := player.IncomingAttacks[0]; len(player.IncomingAttacks) > 0; player.IncomingAttacks = player.IncomingAttacks[1:] {
		dmg := 0
		for name, count := range attack.Units {
			dmg += count * units[name].Attack
		}
		fmt.Printf("Player '%d' is attacked with %d damage...\n", player.ID, dmg)
		player.BaseHealth -= dmg
	}
	return nil
}
