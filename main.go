package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/doc97/zrush/zdata"
	"github.com/doc97/zrush/zsteps"
)

func printBanner() {
	fmt.Println("=== ZRush ===")
	fmt.Println()
	fmt.Println("Creator: Daniel Riissanen (@doc97)")
	fmt.Println("Year: 2019")
	fmt.Println()
	fmt.Println("-------------")
	fmt.Println()
}

func initialize() ([]zdata.PlayerData, map[string]zdata.UnitData, error) {
	fmt.Print("Number of players? ")

	var playerCount int
	if _, err := fmt.Scan(&playerCount); err != nil {
		return nil, nil, fmt.Errorf("could not read input: %s", err)
	}
	if playerCount < 2 {
		return nil, nil, fmt.Errorf("cannot start a game with less than 2 players")
	}

	var players []zdata.PlayerData
	for id := 1; id <= playerCount; id++ {
		players = append(players, zdata.PlayerData{
			ID:         id,
			BaseHealth: 20,
			Minerals:   0,
			Gas:        0,
			Units: map[string]int{
				zdata.Drone:        2,
				zdata.Zergling:     0,
				zdata.Hydralisk:    0,
				zdata.Mutalisk:     0,
				zdata.SporeCrawler: 0,
			},
		})
	}

	units := map[string]zdata.UnitData{
		zdata.Drone:        zdata.UnitData{Name: "Drone", Attack: 0, Defense: 3, Minerals: 1, Gas: 0},
		zdata.Zergling:     zdata.UnitData{Name: "Zergling", Attack: 2, Defense: 1, Minerals: 1, Gas: 0},
		zdata.Hydralisk:    zdata.UnitData{Name: "Hydralisk", Attack: 2, Defense: 3, Minerals: 3, Gas: 1},
		zdata.Mutalisk:     zdata.UnitData{Name: "Mutalisk", Attack: 2, Defense: 1, Minerals: 2, Gas: 2},
		zdata.SporeCrawler: zdata.UnitData{Name: "Spore Crawler", Attack: 0, Defense: 1, Minerals: 1, Gas: 1},
	}
	return players, units, nil
}

func run(players []zdata.PlayerData, units map[string]zdata.UnitData) (int, error) {
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
		fmt.Printf("Base health: %d\n\n", player.BaseHealth)

		if err := zsteps.StepResources(player); err != nil {
			return -1, err
		}

		if err := zsteps.StepEvolve(player); err != nil {
			return -1, err
		}

		if err := zsteps.StepMorph(player, units); err != nil {
			return -1, err
		}

		if err := zsteps.StepDefend(players, units); err != nil {
			return -1, err
		}

		err := zsteps.StepAttack(player, len(players), units)
		if err != nil {
			return -1, err
		}

		pIdx = (pIdx + 1) % len(players)
	}
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
