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
	ID         int
	baseHealth int
	minerals   int
	gas        int
	units      map[string]int
}

// AttackData holds information about an incoming attack.
type AttackData struct {
	defenderID int
	attackerID int
	units      map[string]int
}

// UnitData holds information about a unit.
type UnitData struct {
	name     string
	attack   int
	defense  int
	minerals int
	gas      int
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
		players = append(players, PlayerData{
			ID:         id,
			baseHealth: 20,
			minerals:   0,
			gas:        0,
			units: map[string]int{
				Drone:        2,
				Zergling:     0,
				Hydralisk:    0,
				Mutalisk:     0,
				SporeCrawler: 0,
			},
		})
	}

	units := map[string]UnitData{
		Drone:        UnitData{name: "Drone", attack: 0, defense: 3, minerals: 1, gas: 0},
		Zergling:     UnitData{name: "Zergling", attack: 2, defense: 1, minerals: 1, gas: 0},
		Hydralisk:    UnitData{name: "Hydralisk", attack: 2, defense: 3, minerals: 3, gas: 1},
		Mutalisk:     UnitData{name: "Mutalisk", attack: 2, defense: 1, minerals: 2, gas: 2},
		SporeCrawler: UnitData{name: "Spore Crawler", attack: 0, defense: 1, minerals: 1, gas: 1},
	}
	return players, units, nil
}

func run(players []PlayerData, units map[string]UnitData) (int, error) {
	i := 0
	unitKeys := make([]string, len(units))
	for k := range units {
		unitKeys[i] = k
		i++
	}

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

		attack, err := stepAttack(player, len(players), units)
		if err != nil {
			return -1, err
		}

		if attack != nil {
			dmg := 0
			for name, count := range attack.units {
				dmg += count * units[name].attack
			}
			fmt.Printf("Attacking player '%d' with %d damage...\n", attack.defenderID, dmg)
			players[attack.defenderID-1].baseHealth -= dmg
		}

		pIdx = (pIdx + 1) % len(players)
	}
}

func stepResources(player *PlayerData) error {
	slots, err := subStepSelectResourceSlots(player.units[Drone])
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
		// Exit condition
		if !subStepPrintAffordableUnits(player, units) {
			fmt.Printf("Morphing units...\n\n")
			return nil
		}

		fmt.Printf("\nMorph a unit, 'x' to skip [d,z,h,m,s]? ")
		morphUnits, str, err := readUnitString()
		if err != nil {
			return err
		}
		if str == "x" {
			return nil
		}

		costMinerals, costGas := 0, 0
		for name, count := range morphUnits {
			costMinerals += count * units[name].minerals
			costGas += count * units[name].gas
		}

		if costMinerals > player.minerals || costGas > player.gas {
			fmt.Println("Insufficient funds!")
			continue
		}

		for name, count := range morphUnits {
			player.units[name] += count
		}

		player.minerals -= costMinerals
		player.gas -= costGas
	}
}

func stepAttack(player *PlayerData, playerCount int, units map[string]UnitData) (*AttackData, error) {
	if player.units[Zergling] == 0 && player.units[Hydralisk] == 0 && player.units[Mutalisk] == 0 {
		return nil, nil
	}

	subStepPrintAttackUnits(player, units)

	attackUnits, err := subStepSelectAttackUnits(player)
	if err != nil {
		return nil, fmt.Errorf("failed to select units: %v", err)
	}

	targetPlayer := subStepSelectAttackTarget(playerCount, player.ID)
	attack := AttackData{
		attackerID: player.ID,
		defenderID: targetPlayer,
		units:      attackUnits,
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

func subStepPrintAttackUnits(player *PlayerData, units map[string]UnitData) {
	fmt.Println("Available offensive units:")

	nonAttackUnits := map[string]struct{}{Drone: {}, SporeCrawler: {}}
	for key, count := range player.units {
		if _, contains := nonAttackUnits[key]; contains {
			continue
		}

		if count == 1 {
			fmt.Printf("\t%d %s\n", count, units[key].name)
		} else if count > 1 {
			fmt.Printf("\t%d %ss\n", count, units[key].name)
		}
	}
}

func subStepSelectAttackUnits(player *PlayerData) (map[string]int, error) {
	fmt.Printf("\nAttack with (format: 'zzh')? ")
	units, _, err := readUnitString()
	if err != nil {
		return nil, err
	}
	for name := range units {
		units[name] = zutil.Min(units[name], player.units[name])
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

func readUnitString() (map[string]int, string, error) {
	var unitStr string
	if _, err := fmt.Scan(&unitStr); err != nil {
		return nil, "", fmt.Errorf("could not read input: %v", err)
	}

	units := map[string]int{
		Drone:        0,
		Zergling:     0,
		Hydralisk:    0,
		Mutalisk:     0,
		SporeCrawler: 0,
	}
	for _, unitRune := range strings.ToLower(unitStr) {
		switch unitRune {
		case 'd':
			units[Drone]++
		case 'z':
			units[Zergling]++
		case 'h':
			units[Hydralisk]++
		case 'm':
			units[Mutalisk]++
		case 's':
			units[SporeCrawler]++
		}
	}

	return units, unitStr, nil
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
