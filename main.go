package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/doc97/zrush/zutil"
)

// PlayerData holds information about a player state
// like resources and units.
type PlayerData struct {
	ID            int
	baseHealth    int
	minerals      int
	gas           int
	drones        int
	zerglings     int
	hydralisks    int
	mutalisks     int
	sporeCrawlers int
}

// AttackData holds information about an incoming attack.
type AttackData struct {
	defenderID int
	attackerID int
	zerglings  int
	hydralisks int
	mutalisks  int
}

// UnitData holds information about a unit.
type UnitData struct {
	name     string
	attack   int
	defense  int
	minerals int
	gas      int
}

func printBanner() {
	fmt.Println("=== ZRush ===")
	fmt.Println()
	fmt.Println("Creator: Daniel Riissanen (@doc97)")
	fmt.Println("Year: 2019")
	fmt.Println()
	fmt.Println("-------------")
	fmt.Println()
}

func initialize() ([]PlayerData, map[string]UnitData, error) {
	fmt.Print("Number of players? ")

	var playerCount int
	if _, err := fmt.Scan(&playerCount); err != nil {
		return nil, nil, fmt.Errorf("could not read input: %s", err)
	}
	if playerCount < 2 {
		return nil, nil, fmt.Errorf("cannot start a game with less than 2 players")
	}

	var players []PlayerData
	for id := 1; id <= playerCount; id++ {
		players = append(players, PlayerData{id, 20, 0, 0, 2, 0, 0, 0, 0})
	}

	units := map[string]UnitData{
		"Drone":         UnitData{name: "Drone", attack: 0, defense: 3, minerals: 1, gas: 0},
		"Zergling":      UnitData{name: "Zergling", attack: 2, defense: 1, minerals: 1, gas: 0},
		"Hydralisk":     UnitData{name: "Hydralisk", attack: 2, defense: 3, minerals: 3, gas: 1},
		"Mutalisk":      UnitData{name: "Mutalisk", attack: 2, defense: 1, minerals: 2, gas: 2},
		"Spore Crawler": UnitData{name: "Spore Crawler", attack: 0, defense: 1, minerals: 1, gas: 1},
	}
	return players, units, nil
}

func run(players []PlayerData, units map[string]UnitData) (int, error) {
	pIdx, round := 0, 0
	for {
		if pIdx == 0 {
			round++
			fmt.Printf("\nROUND %d\n----\n", round)
		}
		player := &players[pIdx]

		fmt.Printf("\nPLAYER %d TURN\n", player.ID)
		fmt.Printf("Base health: %d\n\n", player.baseHealth)

		if err := stepResources(player); err != nil {
			return -1, err
		}

		stepEvolve(player)

		if err := stepMorph(player, units); err != nil {
			return -1, err
		}

		attack, err := stepAttack(player, len(players))
		if err != nil {
			return -1, err
		}

		if attack != nil {
			dmg := attack.zerglings*units["Zergling"].attack + attack.hydralisks*units["Hydralisk"].attack + attack.mutalisks*units["Mutalisk"].attack
			fmt.Printf("Attacking player '%d' with %d damage...\n", attack.defenderID, dmg)
			players[attack.defenderID-1].baseHealth -= dmg
		}

		pIdx = (pIdx + 1) % len(players)
	}
}

func stepResources(player *PlayerData) error {
	slots, err := subStepSelectResourceSlots(player.drones)
	if err != nil {
		return err
	}

	min, gas := subStepGenerateResources(slots)
	player.minerals += min
	player.gas += gas
	return nil
}

func stepEvolve(player *PlayerData) {

}

func stepMorph(player *PlayerData, units map[string]UnitData) error {
	for {
		if !subStepPrintAffordableUnits(player, units) {
			return nil
		}

		fmt.Printf("\nMorph a unit, 'x' to skip [d,z,h,m,s]? ")
		d, z, h, m, s, str, err := readUnitString()
		if err != nil {
			return err
		}
		if str == "x" {
			continue
		}

		costMinerals := d*units["Drone"].minerals +
			z*units["Zergling"].minerals +
			h*units["Hydralisk"].minerals +
			m*units["Mutalisk"].minerals +
			s*units["Spore Crawler"].minerals
		costGas := d*units["Drone"].gas +
			z*units["Zergling"].gas +
			h*units["Hydralisk"].gas +
			m*units["Mutalisk"].gas +
			s*units["Spore Crawler"].gas

		if costMinerals > player.minerals || costGas > player.gas {
			fmt.Println("Insufficient funds!")
			continue
		}

		player.drones += d
		player.zerglings += z
		player.hydralisks += h
		player.mutalisks += m
		player.sporeCrawlers += s

		player.minerals -= costMinerals
		player.gas -= costGas
		fmt.Printf("Morphing units...\n")
		fmt.Printf("DEBUG: Morph(%d, %d, %d, %d, %d)\n", d, z, h, m, s)
		fmt.Printf("DEBUG: Total(%d, %d, %d, %d, %d)\n", player.drones, player.zerglings, player.hydralisks, player.mutalisks, player.sporeCrawlers)

		fmt.Println()
		return nil
	}
}

func stepAttack(player *PlayerData, playerCount int) (*AttackData, error) {
	if player.zerglings == 0 && player.hydralisks == 0 && player.mutalisks == 0 {
		return nil, nil
	}

	subStepPrintAttackUnits(player)

	z, h, m, err := subStepSelectAttackUnits(player)
	if err != nil {
		return nil, fmt.Errorf("failed to select units: %v", err)
	}

	targetPlayer := subStepSelectAttackTarget(playerCount, player.ID)
	attack := AttackData{
		attackerID: player.ID,
		defenderID: targetPlayer,
		zerglings:  z,
		hydralisks: h,
		mutalisks:  m,
	}

	fmt.Println()
	return &attack, nil
}

func subStepSelectResourceSlots(drones int) ([]int, error) {
	fmt.Printf("%d drones to use.\n\n", drones)

	var slots []int
	for i := drones; i > 0; {
		fmt.Print("Select resource slot [1 (mineral), 2 (mineral), 3 (vespene gas)]? ")

		var slot int
		if _, err := fmt.Scan(&slot); err != nil {
			return slots, fmt.Errorf("could not read input: %v", err)
		}
		if slot < 1 || slot > 3 {
			fmt.Println("invalid resource slot!")
			continue
		}

		slots = append(slots, slot)
		i--
	}
	return slots, nil
}

func subStepGenerateResources(slots []int) (int, int) {
	min, gas := 0, 0
	dice := rand.Intn(6) + 1
	fmt.Printf("DEBUG: dice was %d\n", dice)

	for i := 1; i <= 3; i++ {
		if dice <= 3 && i != dice {
			continue
		}

		count := 0
		for _, slot := range slots {
			if slot == i {
				count++
			}
		}

		if i == 3 && (dice == 3 || dice == 5 || dice == 6) {
			gas += count
			if count > 1 {
				gas++
			}
		} else if i == 2 && (dice == 2 || dice == 4 || dice == 6) {
			min += count
			if count > 1 {
				min++
			}
		} else if i == 1 && (dice == 1 || dice == 4 || dice == 5) {
			min += count
			if count > 1 {
				min++
			}
		}
	}

	fmt.Println()
	return min, gas
}

func subStepPrintAffordableUnits(player *PlayerData, units map[string]UnitData) bool {
	affordableUnits := make(map[string]UnitData)
	for name, unit := range units {
		if player.minerals >= unit.minerals && player.gas >= unit.gas {
			affordableUnits[name] = unit
		}
	}

	if len(affordableUnits) == 0 {
		return false
	}

	fmt.Printf("Minerals: %d\nVespene Gas: %d\n\n", player.minerals, player.gas)
	fmt.Println("You can morph the following units:")
	for _, unit := range affordableUnits {
		fmt.Printf("%s (%d minerals, %d gas)\n", unit.name, unit.minerals, unit.gas)
	}

	return true
}

func subStepPrintAttackUnits(player *PlayerData) {
	fmt.Println("Available offensive units:")

	if player.zerglings == 1 {
		fmt.Printf("\t%d zergling\n", player.zerglings)
	} else if player.zerglings > 1 {
		fmt.Printf("\t%d zerglings\n", player.zerglings)
	}

	if player.hydralisks == 1 {
		fmt.Printf("\t%d hydralisk\n", player.hydralisks)
	} else if player.hydralisks > 1 {
		fmt.Printf("\t%d hydralisks\n", player.hydralisks)
	}

	if player.mutalisks == 1 {
		fmt.Printf("\t%d mutalisk\n", player.mutalisks)
	} else if player.mutalisks > 1 {
		fmt.Printf("\t%d mutalisks\n", player.mutalisks)
	}
}

func subStepSelectAttackUnits(player *PlayerData) (int, int, int, error) {
	fmt.Printf("\nAttack with (format: 'zzh')? ")
	_, z, h, m, _, _, err := readUnitString()
	if err != nil {
		return 0, 0, 0, err
	}
	z = zutil.Min(z, player.zerglings)
	h = zutil.Min(h, player.hydralisks)
	m = zutil.Min(m, player.mutalisks)
	return z, h, m, nil
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

func readUnitString() (int, int, int, int, int, string, error) {
	var unitStr string
	if _, err := fmt.Scan(&unitStr); err != nil {
		return 0, 0, 0, 0, 0, "", fmt.Errorf("could not read input: %v", err)
	}

	d, z, h, m, s := 0, 0, 0, 0, 0
	for _, unitRune := range strings.ToLower(unitStr) {
		switch unitRune {
		case 'd':
			d++
		case 'z':
			z++
		case 'h':
			h++
		case 'm':
			m++
		case 's':
			s++
		}
	}

	return d, z, h, m, s, unitStr, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	printBanner()

	players, units, err := initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	winner, err := run(players, units)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	fmt.Println("Winner: ", winner)
}
