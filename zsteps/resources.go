package zsteps

import (
	"fmt"
	"math/rand"

	"github.com/doc97/zrush/zdata"
)

// StepResources performs the resource step
func StepResources(player *zdata.PlayerData) error {
	slots, err := selectResourceSlots(player.Units[zdata.Drone])
	if err != nil {
		return err
	}

	min, gas := generateResources(slots)
	player.Minerals += min
	player.Gas += gas
	return nil
}

func selectResourceSlots(drones int) ([]int, error) {
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

func generateResources(slots []int) (int, int) {
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
